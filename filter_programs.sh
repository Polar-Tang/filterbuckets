#!/bin/bash

input_file="intigriti2.txt"

# Process the file line by line
while IFS= read -r line; do
    # Remove leading/trailing spaces and check conditions
    trimmed_line="$(echo "$line" | awk '{$1=$1};1')"

    # Skip empty lines or lines matching specific patterns
    if [[ -z "$trimmed_line" ]] || [[ "$trimmed_line" == "View program" ]] || [[ "$trimmed_line" == "Responsible Disclosure" ]] || [[ ${#trimmed_line} -gt 25 ]]; then
        continue
    fi

    # Print valid lines
    printf "%s\n" "$trimmed_line"
done < "$input_file"
