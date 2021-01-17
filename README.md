![shopify-mono-black](https://user-images.githubusercontent.com/46195831/104846062-e601f400-58d8-11eb-9110-62aff1ed97e8.png)
# **SHOPIFY CHALLENGE**

This is my submission for the Shopify 2021 Challenge. I have created an Image Repository written in Golang. This is a basic API with three endpoints for user signup, Login & Image Upload.
In Building this project, I used: 
 - Gorilla Mux as the HTTP Router
 - MongoDB as the database
 - Google Cloud Storage as a file storage service

# **Try It Out**
To try out this project, You are required to have a Golang installed on your computer, A MongoDB database & a GCP account -- 

 - download your GCP service account file, Name it ***keys.json*** and store it in the root of this project directory. 
 - A MongoDB URI is required, add it to the ***.env*** file.
 - Create a Cloud storage bucket and add the name to the ***.env*** file.
 -  Execute go run main.go
