#!/bin/bash

# Output file to save all the extracted program URLs
output_file="firebounty_programs.txt"
> "$output_file"  # Clear the file if it already exists

# Base URL
base_url="https://firebounty.com/?page="

# Number of pages to scrape (adjust as needed)
max_pages=645  # You can change this number if more pages exist

# Timeout settings
wget_timeout=10  # Timeout in seconds
wget_retries=2   # Number of retry attempts

# Loop through all the pages
for ((i=400; i<=max_pages; i++)); do
    echo "Processing page $i..."

    # Temporary file to save the downloaded HTML content
    tmp_file="index.html?page=$i"

    # Download the HTML content with timeout and retries
    wget --timeout="$wget_timeout" --tries="$wget_retries" -q -O "$tmp_file" "${base_url}${i}"

    # Check if the file is not empty
    if [[ -s "$tmp_file" ]]; then
        # Apply the provided one-liner to extract program URLs
        cat "$tmp_file" | grep "<div class='box'>" -C 2 | grep "center-helper" | \
            awk '{print $6}' | sed -n "s/.*='\([^']*\)'.*/\1/p" >> "$output_file"

        echo "Page $i processed successfully."
    else
        echo "Page $i failed to download or is empty. Skipping."
    fi

    # Optional: Clean up the temporary file
    rm -f "$tmp_file"
done

echo "All pages processed. Results saved to $output_file."

