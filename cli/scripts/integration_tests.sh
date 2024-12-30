#!/bin/sh

files="./scripts/test1.txt ./scripts/test2.txt ./scripts/test3.txt ./scripts/test4.txt"

shell_output="shell_output.txt"
bash_output="bash_output.txt"

for file in $files; do
    echo "Testing with file: $file"

    ./shell < "$file" > "$shell_output"
    bash < "$file" > "$bash_output"

    if diff "$shell_output" "$bash_output" > /dev/null; then
        echo "[PASS] Outputs match for file: $file"
    else
        echo "[FAIL] Outputs differ for file: $file"
        echo "Differences:"
        diff "$shell_output" "$bash_output"
    fi

done

rm -f "$shell_output" "$bash_output"