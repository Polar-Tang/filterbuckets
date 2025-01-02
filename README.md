#### **1. Install Go:**
   - **Download:** Go to [golang.org/dl](https://golang.org/dl/) and download the installer for your OS.
   - **Install:** Follow the instructions for your operating system.

#### **2. Verify Installation:**
   Open your terminal and run:
   ```bash
   go version
   ```
   You should see the installed Go version (e.g., `go version go1.21.0`).

---

#### **3. Set your api key:**
   Set the greyhat api that you will find in https://grayhatwarfare.com/account/settings:
   ```bash
   echo "API_KEY" > ./sessionCookie
   ```

---

#### **4. Install the tool:**

   ```sh
   go install github.com/Polar-Tang/filterbuckets@v0.2.0
   echo 'alias filterbuckets="~/go/bin/filterbuckets"' >> ~/.profile
   source ~/.profile
   ```
 
---

#### **5. Run the Program:**
   ```bash
   go run main.go
   ```

---
