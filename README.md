# books_club
A simple CRUD app in golang using net/http package for books enthusiasts.

## APIs:

### GET:
- `/books`  
  returns all the books from the database
- `/author/{id}`  
  returns the detail of author 

### POST:
- `/book`  
  creates an entry of book in the database
- `/author`  
  creates an entry of author in the database
- `/rate-book`  
  rates a book and returns the updated average rating
