package components

/*
	Create a constant for every single error
*/

const InternalServerError string = "{\"code\": 500, \"message\": \"Internal Server Error\"}"
const BadRequestError string = "{\"code\": 400, \"message\": \"Bad Request\"}"
const UnauthorizedError string = "{\"code\": 401, \"message\": \"Unauthorized\"}"
const NotFoundError string = "{\"code\": 404, \"message\": \"Not Found\"}"
const ConflictError string = "{\"code\": 409, \"message\": \"Conflict\"}"
const ForbiddenError string = "{\"code\": 403, \"message\": \"Forbidden\"}"

const NoContent string = "{\"code\": 204, \"message\": \"No Content\"}"
