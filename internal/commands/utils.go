package commands

import "github.com/jessevdk/go-flags"

// arg_parse парсит аргументы команды в переданную структуру
func arg_parse[Rcv any, PtrRcv *Rcv](rcv PtrRcv, args []string) error {
	parser_options := flags.Options(flags.PrintErrors | flags.IgnoreUnknown)
	parser := flags.NewParser(rcv, parser_options)

	_, err := parser.ParseArgs(args)
	if err != nil {
		return err
	}
	return nil
}
