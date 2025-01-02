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

   - **Linux:** 
     ```bash
     sudo apt-get install tesseract-ocr libtesseract-dev 
     ```
   - **Mac:** 
     ```bash
     brew install tesseract
     ```
   - **Windows:** 
     Download the installer from [Tesseract GitHub](https://github.com/tesseract-ocr/tesseract/releases/download/5.5.0/tesseract-ocr-w64-setup-5.5.0.20241111.exe) and follow the setup instructions.

---

#### **5. Run the Program:**
   ```bash
   go run main.go
   ```

---
