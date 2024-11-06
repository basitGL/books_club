# books_club
A simple CRUD app in golang using net/http package for books enthusiasts.

## APIs:

### Public Routes:

Routes that are public and don't require user credentials

#### POST:
- `/register`
  creates a user, saves record in mysql db, and returns the user token
- `/login`
  returns the user token if successful else returns appropriate message


### Protected Routes:

Routes that protected through a middleware of user token

#### GET:
- `/books`  
  returns all the books from the database
- `/author/{id}`  
  returns the detail of author 

#### POST:
- `/book`  
  creates an entry of book in the database
- `/author`  
  creates an entry of author in the database
- `/rate-book`  
  rates a book and returns the updated average rating
