# Airbnb mockup project
Project made for Service Oriented Architecture

## Short synopsis

Airbnb clone:
- A microservice web application that simulates Airbnb
  
- There are 3 types of users: Host, Guest and Unauthenticated User
- You can register as either a Guest or a Host
- Through a RBAC the application gives you permission to use different functionalities of the application.

- Guests have the option to view accommodations plan their date of visiting and make a reservation if the desired accommodation is available for the defined period.
- All reserved accommodations are displayed on the Guests profile and he can manage them accordingly.
- Each reservation can be canceled before the defined period starts.
- At the end of each reservation Guests are given an option to review both the accommodation and the Host.
- The ratings contain a 1-5 star option and comment given gy the Guest.
- All ratings are displayed on the individual page of the accommodation that was reviewed as well as on the profile of both Guest and Host.
- Hosts have the option of making a accommodation, providing a availability date, as well as managing CRUD operations on that accommodation.


## Tools used for the project
- Back end: Go
- Front end: Angular
- Database: MongoDB, CassandraDB
- Tools: Docker, Redis (To be implemented), Kubernetes (To be implemented)
