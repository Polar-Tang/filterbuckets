#!/bin/bash

input_file="intigriti2.txt"

while IFS= read -r line; do
    trimmed_line="$(echo "$line" | awk '{$1=$1};1')"

    if [[ -z "$trimmed_line" ]] || [[ "$trimmed_line" == "View program" ]] || [[ "$trimmed_line" == "Responsible Disclosure" ]] || [[ ${#trimmed_line} -gt 25 ]]; then
        continue
    fi

    printf "%s\n" "$trimmed_line"
done < "$input_file"
