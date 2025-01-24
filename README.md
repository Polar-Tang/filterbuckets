### Description

Filterbuckets utilizes greyhat api to download public buckets and compile regex to filter for sensitive data

### Installation

   ```sh
   go install github.com/Polar-Tang/filterbuckets@v0.2.2
   echo 'alias filterbuckets="~/go/bin/filterbuckets"' >> ~/.profile
   source ~/.profile
   ```
 
---

### Usage
Four flags, x flag is the only that is always required, this should be the filename where the extension is the key and the value it's an array with the regex to search inside every file.

```
Usage of /home/pull/go/bin/filterbuckets:
  -b string
        Bucket name.
  -c int
        Concurrency limit for processing (default: 3). (default 3)
  -w string
        Path to the file containing keywords.
  -x string
        Path to the file containing extensions.

```

#### Set your api key:
   Set the greyhat api that you will find in https://grayhatwarfare.com/account/settings:
   ```bash
   echo "API_KEY" > ./sessionCookie
   ```

---
### Motivation
If you found your self using greyhat without any tool probably you'd get frustrate, there are unlimited files with random data, i personally were doing a tedious routine, such as:
- enter the link,
- see data, 
- leave, 
- repeate  
Filterbuckets do everything at once, 