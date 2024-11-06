# Online Learning Platform - Golang

## Introduction

The **Online Learning Platform** project is built using **Golang**, aimed at providing a flexible and user-friendly online learning platform for both students and instructors. This application allows instructors to create courses, lectures, and materials, helping students participate in learning and track their progress.

## Key Features

- **Course Management**: Create, update, and delete online courses.
- **Lecture Management**: Add lectures to courses, manage content, and related files.
- **Course Enrollment**: Users can enroll in available courses.
- **File Upload**: Use Cloudinary to upload and manage lecture materials.

## Technologies Used

- **Golang**: The main programming language for the backend.
- **Cloudinary**: Service for storing images and videos.
- **AWS S3**: Service for storing PDF files.
- **MySQL**: Database for storing course, lecture, and user information.
- **Docker**: For deploying the application and ensuring a consistent development environment.

## Installation

### System Requirements

- **Go**: Version 1.22
- **MySQL**: Database for running the application.
- **Docker**: To run the application in a container.

### Installation Guide

1. **Clone the repository**

   ```bash
   git clone https://github.com/thanhminhdev2000/online-learning-golang.git
   cd online-learning-golang
   ```

2. **Configure the environment**
   Create a `.env` file with the following content:

   ```env
   DB_CONNECTION=user:user_pw@tcp(localhost:3306)/online-learning
   PORT=
   API_PREFIX=
   CLIENT_URL=
   JWT_KEY=
   SMTP_HOST=smtp.gmail.com
   SMTP_EMAIL=
   SMTP_PASSWORD=
   CLOUDINARY_CLOUD_NAME=
   CLOUDINARY_API_KEY=
   CLOUDINARY_API_SECRET=
   AWS_ACCESS_KEY_ID=
   AWS_SECRET_ACCESS_KEY=
   AWS_S3_BUCKET_NAME=
   AWS_REGION=us-east-1
   AWS_STORAGE=
   CLOUDINARY_STORAGE=
   COOKIE_DOMAIN=localhost
   ENV=development
   ```

3. **Run the application using Docker**

   ```bash
   docker-compose up --build -d
   ```

   ```bash
   docker-compose up -d
   ```

4. **Access the application**
   The application will run at: `http://localhost:8080`

## Usage

For testing purposes, you can use the following accounts:

**Admin Account**:

- Username: admin1
- Password: zzzxxx

**User Account**:

- Username: user11
- Password: zzzxxx

## Folder Structure

- **/controllers**: Contains controllers for handling HTTP requests.
- **/models**: Contains models for course, lecture, and user data.
- **/routes**: Defines the system's API endpoints.
- **/utils**: Contains utilities such as database connections and file uploads.

## Contribution

If you want to contribute to the project, please create a **Pull Request** or a new **Issue**. All contributions are welcome!

## Contact

- **Author**: Nguyen Thanh Minh
- **Email**: [thanhminh.nguyendev@gmail.com](mailto:thanhminh.nguyendev@gmail.com)
- **Deployed at**: [http://52.90.82.84/](http://52.90.82.84/)

Thank you for your interest in our project!
