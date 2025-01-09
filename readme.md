# Personal Blog

This project is a simple personal blog that allows the admin to add, edit, and delete articles through a web interface. It uses Go for the backend and stores articles in JSON format.

## Getting Started

To run this project on your local machine, follow these steps.

### 1. Clone the Repository

Clone the repository using Git:

```bash
git clone https://github.com/demenkoeugene/personal_blog.git
cd personal_blog
```

### 2. Set Up the config.yaml File

Create a config.yaml file in the /config directory of your project to store the configuration settings for your blog. Example:
```yaml
auth:
  username: "admin"
  password: "yourpassword"
```

This file allows you to configure 	admin: Admin credentials for logging into the admin interface.

### 3. Run the Server

To start the server, use the following command:
```bash
go run main.go
```

The server will start on the port 8080.

### 4. Create an Article via the Admin Interface

1. Go to http://localhost:8080/admin.
2.	Log in using the admin credentials from config.yaml.
3.	To create a new article, click “Add New Article”.
4.	Fill in the title, content, and publication date for the article.
5.	After completing the form, click “Submit” to save the article.

###  5. Editing and Deleting Articles
* Edit an Article: To edit an article, go to the edit page via a URL like http://localhost:8080/edit/{id}.
* Delete an Article: To delete an article, go to the delete page via a URL like http://localhost:8080/delete/{id}.

### 6. Technical Details
* Articles are stored in JSON format, with each article saved in a separate file.
* The admin has the ability to create, edit, and delete articles via the admin interface.


