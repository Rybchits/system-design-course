package components

import (
	"math"
	"roguelike/packages/ecs"
	"time"
)

type MobBehavior struct {
	Strategy MobStrategy
}

type MobStrategy interface {
	Act(entity *ecs.Entity, em ecs.EntityManager)
}

func (a *MobBehavior) Mask() uint64 {
	return MaskMobStrategy
}

func NewMobBehavior(strategy MobStrategy) *MobBehavior {
	return &MobBehavior{Strategy: strategy}
}

type aggressiveEnemyStrategy struct {
	delayMillisec         int
	previousStepTimestamp time.Time
}

func NewAggressiveStrategy(delayMillisec int) *aggressiveEnemyStrategy {
	return &aggressiveEnemyStrategy{
		delayMillisec:         delayMillisec,
		previousStepTimestamp: time.Now(),
	}
}

func (b *aggressiveEnemyStrategy) Act(entity *ecs.Entity, em ecs.EntityManager) {
	// Проверка времени задержки
	if time.Since(b.previousStepTimestamp).Milliseconds() < int64(b.delayMillisec) {
		return
	}

	// Реализация агрессивного поведения: атакуют игрока, как только его видят
	player := em.Get("player")
	location := em.Get("location").Get(MaskLocation).(*Location)
	entityPosOrNil := entity.Get(MaskPosition)

	if player == nil || location == nil || entityPosOrNil == nil {
		return
	}

	entityPos := entityPosOrNil.(*Position)
	playerPos := player.Get(MaskPosition).(*Position)

	// BFS для поиска пути к игроку
	queue := []Position{*entityPos}

	previous := make([][]*Position, location.MapSize.Height)
	for i := range previous {
		previous[i] = make([]*Position, location.MapSize.Width)
	}
	previous[entityPos.Y][entityPos.X] = nil

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		adjacentPositions := []Position{
			{X: current.X + 1, Y: current.Y},
			{X: current.X - 1, Y: current.Y},
			{X: current.X, Y: current.Y + 1},
			{X: current.X, Y: current.Y - 1},
		}

		for _, pos := range adjacentPositions {
			if pos.X == (*playerPos).X && pos.Y == (*playerPos).Y {
				previous[pos.Y][pos.X] = &current
				queue = nil
				break
			}

			if !location.IsAvailablePosition(pos) || !pos.IsFree(em) || previous[pos.Y][pos.X] != nil {
				continue
			}
			queue = append(queue, pos)
			previous[pos.Y][pos.X] = &current
		}
	}

	if previous[(*playerPos).Y][(*playerPos).X] != nil {
		path := []Position{}
		for at := playerPos; at != nil; at = previous[(*at).Y][(*at).X] {
			path = append(path, *at)
		}

		if len(path) > 1 {
			nextPos := path[len(path)-2]
			entity.Add(NewMovement().WithNext(nextPos))
		}
	}

	// Обновляем время последнего шага
	b.previousStepTimestamp = time.Now()
}

type passiveStrategy struct{}

func (b *passiveStrategy) Act(entity *ecs.Entity, em ecs.EntityManager) {
	// Реализация пассивного поведения: просто стоят на месте
}

func NewPassiveStrategy() *passiveStrategy {
	return &passiveStrategy{}
}

type cowardEnemyStrategy struct {
	delayMillisec         int
	previousStepTimestamp time.Time
}

func NewCowardStrategy(delayMillisec int) *cowardEnemyStrategy {
	return &cowardEnemyStrategy{
		delayMillisec:         delayMillisec,
		previousStepTimestamp: time.Now(),
	}
}

func (b *cowardEnemyStrategy) Act(entity *ecs.Entity, em ecs.EntityManager) {
	// Проверка времени задержки
	if time.Since(b.previousStepTimestamp).Milliseconds() < int64(b.delayMillisec) {
		return
	}

	player := em.Get("player")
	location := em.Get("location").Get(MaskLocation).(*Location)
	entityPosOrNil := entity.Get(MaskPosition)

	if player == nil || location == nil || entityPosOrNil == nil {
		return
	}

	entityPos := entityPosOrNil.(*Position)
	playerPos := player.Get(MaskPosition).(*Position)

	// Определяем смежные клетки (влево, вправо, вверх, вниз)
	adjacentPositions := []Position{
		{X: entityPos.X + 1, Y: entityPos.Y},
		{X: entityPos.X - 1, Y: entityPos.Y},
		{X: entityPos.X, Y: entityPos.Y + 1},
		{X: entityPos.X, Y: entityPos.Y - 1},
	}

	// Вычисляем манхэттенское расстояние
	manhattanDistance := func(a, b Position) int {
		return int(math.Abs(float64(a.X-b.X)) + math.Abs(float64(a.Y-b.Y)))
	}

	// Ищем клетку с максимальным расстоянием
	var farthestCell *Position = nil
	maxDistance := manhattanDistance(*playerPos, *entityPos)

	for _, neighbor := range adjacentPositions {
		if !location.IsAvailablePosition(neighbor) || !neighbor.IsFree(em) {
			continue
		}

		distance := manhattanDistance(neighbor, *playerPos)
		if distance > maxDistance {
			maxDistance = distance
			farthestCell = &neighbor
		}
	}

	if farthestCell != nil {
		entity.Add(NewMovement().WithNext(*farthestCell))
	}

	b.previousStepTimestamp = time.Now()
}
