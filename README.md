### Description

Filterbuckets utilizes greyhat api to download public buckets and compile regex to filter for sensitive data

### Installation

   ```sh
   go install github.com/Polar-Tang/filterbuckets@latest
   # For Linux/WSL
   echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.profile
   # next, reload your shell configuration
   source ~/.profile
   ```

   ```
   # For Mac OS
   echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.zshrc
   # next, reload your shell configuration
   source ~/.zshrc
   ```
 
 #### Set your api key:
   Set the greyhat api that you will find in https://grayhatwarfare.com/account/settings:
   ```bash
   echo "API_KEY" > ./sessionCookie
   ```

---

### Usage
Usage of /home/pull/go/bin/filterbuckets:
```
  -b string
        Bucket name.
  -c int
        Concurrency limit for processing (default: 3). (default 3)
  -w string
        Path to the file containing keywords.
  -x string
        Path to the file containing extensions.
```

Example:
```
filterbuckets -x testing.json -w intigriti2.txt -c 1 -b platform-dev.storage.googleapis.com
```

For further information you could see the [tutorial of the tool](https://www.youtube.com/watch?v=dTYtrbLA61s&t=102s)
