#!/bin/bash

output_file="firebounty_programs.txt"
> "$output_file" 

base_url="https://firebounty.com/?page="

max_pages=645  

wget_timeout=10  
wget_retries=2   

for ((i=400; i<=max_pages; i++)); do
    echo "Processing page $i..."

    tmp_file="index.html?page=$i"

    wget --timeout="$wget_timeout" --tries="$wget_retries" -q -O "$tmp_file" "${base_url}${i}"

    if [[ -s "$tmp_file" ]]; then
        cat "$tmp_file" | grep "<div class='box'>" -C 2 | grep "center-helper" | \
            awk '{print $6}' | sed -n "s/.*='\([^']*\)'.*/\1/p" >> "$output_file"

        echo "Page $i processed successfully."
    else
        echo "Page $i failed to download or is empty. Skipping."
    fi

    rm -f "$tmp_file"
done

echo "All pages processed. Results saved to $output_file."

