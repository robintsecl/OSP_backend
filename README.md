# OSP_backend
OSP_backend

# Background
This api server provide API endpoints for user (frontend) to access survey and response using Golang and mongoDB with MVC structure

# Prerequisite
- The application will be run on window 11
- Install Docker desktop with latest verison
- Install Git
- Install mongoDB
    * docker pull mongo:8.0.8
    * docker run --name osp-db -d -p 27017:27017 mongo
- Install mongosh v2.5.0
    * https://www.mongodb.com/try/download/shell
- Install Golang
    * https://go.dev/dl/
    * go1.24.2.windows-amd64.msi
- Reboot host to apply environment variable
- Install any API testing tool such as Postman
- Clone the project from https://github.com/robintsecl/OSP_backend.git
- Switch to main brunch
- Ensure port 9091 and 27017 are not blocked by any firewall or being used by other application
- Confirm container osp-db is running
- Change directory to the location that you have cloned the project
- Run "go run main.go"

# Assumption
- This API server should be run on localhost only
- The caller of API server should be frontend application, not individual user, so frontend should have error checking on empty value, survey.type error before passing the data to these api enpoints

# Code explanation

# main.go
line 32: Create context.

line 33-42: Connect to MongoDB at localhost 27017. If there is connection error, print error log and exit the application.

line 43-56: Connect to the collection called "user", check if there is a document with name = admin. If not, create a document with password, then print the admin password to console.

line 58: Initialize validator, will pass to controller later.

line 60-62: Connect to the collection called "survey".
Create Survey service instance and pass surveycollection and context.
Create Survey controller instance and pass surveyservice, usercollection and context.

line 64-66: Connect to the collection called "response".
Create Response service instance and pass responsecollection, surveycollection and context.
Create Response controller instance and pass responseservice, surveyservice, usercollection and context.

line 68-69: Initialize gin route with release mode.

line 72-79: Main function, run server and append survey and respnse path after basepath. Disconnect Mongo after main exited.

# models
Survey: Survey and Question are created for survey structure.
Response: Response and ResponseAnswer are created for response structure.

Survey is like a google form with type and spec defined by creator, and answerer need to get survey token to answer survey, following the type(format) and spec specified by creator. Only 3 type of questions are accepted including Textbox, Multiple Choice, and Likert scale. Creator must input question.type follows the string in constants.go: "TEXTBOX", "MC" and "LS". Otherwises, an error will be thrown. After creator decided the type(format) of question, they need to provide the spec as an array of options. With textbox type, the spec array can be empty since it receive string anyway. But for Multiple Choice and Likert scale, survey should provide an array of option for answerer (frontend) to select. Details will be in utils section below -> checkFormatAndSpec.

You can checkout the /test/ folder, the response.token should be replaced by the value after createSurvey as value to map response with appropriate survey. 

# Survey APIs

# CreateSurvey (Http.POST)
url: http://localhost:9091/osp/survey/create

This api require "name" and "password" query parameter as well as request body.

In controller layer, firstly, utils.CheckAdmin will check if admin name and password in query parameter are correct. If so, bind JSON from request body to survey. Then validate if survey data violated validate requirement in model type. If not, insert create and update date, and pass to SurveyService.CreateSurvey.

In service layer, recursively generate a 5 characters token and check if any existing token in DB, until a brand new token is generated. After that, inject token to Survey token field. Then do some common checking utils.CommonChecking(survey.Questions), more details in utils section. If no error, insert the data.

# GetSurvey (Http.GET)
url: http://localhost:9091/osp/survey/get

This api require "token" query parameter of survey.

In controller layer, firstly, check "token" query parameter is empty. Then get the survey with token in service layer.

In service layer, find data will token and return result.

# GetAll (Http.GET)
url: http://localhost:9091/osp/survey/getall

This api require "name" and "password" query parameter.

In controller layer, firstly, utils.CheckAdmin will check if admin name and password in query parameter are correct. If so, call SurveyService.GetAll in service layer.

In service layer, for loop find all survey data with cursor and assign data to surveys array. Close the cursor and return surveys.

# UpdateSurvey (Http.PUT)
url: http://localhost:9091/osp/survey/update

This api require "name" and "password" query parameter as well as request body.

In controller layer, firstly, utils.CheckAdmin will check if admin name and password in query parameter are correct. If so, bind JSON from request body to survey. Then validate if survey data violated validate requirement in model type. If not, insert update date, and pass to SurveyService.UpdateSurvey.

In service layer, do some common checking utils.CommonChecking(survey.Questions), more details in utils section. If no error, update the survey.

# DeleteSurvey (Http.DELETE)
url: http://localhost:9091/osp/survey/delete

This api require "name", "password" and "token" query parameter.

In controller layer, firstly, utils.CheckAdmin will check if admin name and password in query parameter are correct. If so, check if token query parameter is empty. If not, go to SurveyService.DeleteSurvey with token passed.

In service layer, delete survey with matched token.

# Response APIs

# CreateResponse (Http.POST)
url: http://localhost:9091/osp/responses/create

This api require request body.

In controller layer, firstly, bind JSON from request body to response. Then validate if response data violated validate requirement in model type. If not, insert create and update date, and pass to ResponseService.CreateResponse.

In service layer, use the token value in responses data to check if any matched survey with token. Token will be used to map the survey that the user is choosing and answering. If survey is found, check answer with utils.ResponseInputChecking(survey.Questions, response.ResponseAnswer), more details in utils section. If no error, insert the data.

# GetAll (Http.GET)
url: http://localhost:9091/osp/responses/getall

This api require "name" and "password" query parameter.

In controller layer, firstly, utils.CheckAdmin will check if admin name and password in query parameter are correct. If so, call ResponseService.GetAll in service layer.

In service layer, for loop find all response data with cursor and assign data to responses array. Close the cursor and return responses.

# GetByToken (Http.GET)
url: http://localhost:9091/osp/responses/getbytoken

This api require "name", "password" and "token" query parameter.

In controller layer, firstly, utils.CheckAdmin will check if admin name and password in query parameter are correct. If so, check if token query parameter is empty. If not, call ResponseService.GetByToken.

In service layer, for loop find all response data matches with "token" with cursor and assign data to responses array. Close the cursor and return responses.

# BatchDeleteResponse (Http.DELETE)
url: http://localhost:9091/osp/responses/batchdelete

This api require "name", "password" and "token" query parameter.

In controller layer, firstly, utils.CheckAdmin will check if admin name and password in query parameter are correct. If so, check if token query parameter is empty. If not, go to SurveyService.BatchDeleteByToken with token passed.

In service layer, delete all responses with matched token.

# Utils.go

# CheckAdmin
This function receive 2 arguments including context and usercollection.

Firstly, it will get name and password query parameters in context, then check if it is empty. If not, query user collection with name and password provided by user. If matched document is found, then login succeed.

# CommonChecking

This function receive 1 argument including questions.

This function will loop through the questions array, and check the followings:

checkIsDuplicateTitle -> If there are any duplicated title, set isDupTitle to true
checkFormatAndSpec -> If spec not match the format, set isWrongFormat to true

return the error as per the isDupTitle and isWrongFormat booleans

# checkIsDuplicateTitle

This function receive 3 arguments including currentValue, titleMap and isDup.

This function checks if the titleMap contains key of current value which is the title. If so, set isDup to true, else append the title of currentValue to the map as key.

# checkFormatAndSpec

This function receive 2 arguments including question and isWrongFormat.

This function checks:
1. In case of textbox, return regardless of the array of spec user has input. There is no restriction.
2. In case of Likert scale, if user provide less than 4 items in spec (options), set isWrongFormat to true.
3. In case of Multiple choice, if user provide less than 2 items in spec (options), set isWrongFormat to true.
4. In case of other format, considered to be unknown type, set isWrongFormat to true.

# ResponseInputChecking

This function receive 2 arguments including questions and answers.

This function is to check if answers in response fulfills the spec in questions. Firstly for loop and store all questions in questionMap with key = question.title, value = question. Then loop through answers array. For each element, check if answer.title exists in questionMap. If so, that means we can check the spec to see if user input right answer.
For textbox -> if length of answer is smaller than 1, return error.
For Likert scale and Multiple choice -> if answerer haven't input the answer in question.spec, return error.

# InsertSurveyDate & InsertResponseDate

These functions receive 2 arguments including response/survey and insertType.

These function checks if insertType is update or create. If update, assign time.Now() to updateDate. If create, assign time.Now() to both createDate and updateDate

# Potential issue as a prodiction-wise server
- Insuffficient data type
If the survey contains many questions, the whole survey document may be so large in size after platform scaled up. Separating question and survey into 2 collections, then store the question ids to survey will be a solution for this issue.
- Simple admin authentication method
In this sample api server, it is designed to have admin login in order to perform some create and update action. However, it has fix user and password, and passing the value from query parameter only. A production-wise application should have an authentication server to control APIs access, also have a user collection, some endpoint for user management in db.
- CORS issue
This sample api server should only allow localhost request but not from other ip. A production-wise api server should have CORS setting configurated and additional coding pointing to the domain of frontend.
- Logging
This sample api server only printed log to stdout but not stored in file. A production-wise api server should have logging library like log4j in Java to handle logging and file rotation.
- High availability
The sample api server has one node only. A production-wise api server should have more than one node and allow to failover upon server failure. Also apply load balancer to allocate triffic.
- Security
Reverse Proxy should be applied for load balancing and protecting api-server from attack. Also allow resource forwarding especially when there are several applications under a domain. Furthermore, it handles SSL encryption and decryption to provide a more secure environment.