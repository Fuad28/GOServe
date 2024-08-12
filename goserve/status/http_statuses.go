package status

import (
	"strconv"
)

type HTTPStatus struct {
	Code        int
	Message     string
	Description string
}

type HTTPStatusMap map[int]HTTPStatus

// HTTP Status map used in constructing response string.
var HTTPStatuses = HTTPStatusMap{
	100: {
		Code:        100,
		Message:     "Continue",
		Description: "The server has received the request headers, and the client should proceed to send the request body.",
	},
	101: {
		Code:        101,
		Message:     "Switching Protocols",
		Description: "The requester has asked the server to switch protocols.",
	},
	102: {
		Code:        102,
		Message:     "Processing",
		Description: "This code indicates that the server has received and is processing the request, but no response is available yet. This prevents the client from timing out and assuming the request was lost.",
	},
	103: {
		Code:        103,
		Message:     "Early Hints",
		Description: "Used to return some response headers before final HTTP message.",
	},
	200: {
		Code:        200,
		Message:     "OK",
		Description: "The request is OK (this is the standard response for successful HTTP requests).",
	},
	201: {
		Code:        201,
		Message:     "Created",
		Description: "The request has been fulfilled, and a new resource is created.",
	},
	202: {
		Code:        202,
		Message:     "Accepted",
		Description: "The request has been accepted for processing, but the processing has not been completed.",
	},
	203: {
		Code:        203,
		Message:     "Non-Authoritative Information",
		Description: "The request has been successfully processed, but is returning information that may be from another source.",
	},
	204: {
		Code:        204,
		Message:     "No Content",
		Description: "The request has been successfully processed, but is not returning any content.",
	},
	205: {
		Code:        205,
		Message:     "Reset Content",
		Description: "The request has been successfully processed, but is not returning any content, and requires that the requester reset the document view.",
	},
	206: {
		Code:        206,
		Message:     "Partial Content",
		Description: "The server is delivering only part of the resource due to a range header sent by the client.",
	},
	207: {
		Code:        207,
		Message:     "Multi-Status",
		Description: "The message body that follows is by default an XML message and can contain a number of separate response codes, depending on how many sub-requests were made.",
	},
	208: {
		Code:        208,
		Message:     "Already Reported",
		Description: "The members of a DAV binding have already been enumerated in a preceding part of the (multistatus) response, and are not being included again.",
	},
	218: {
		Code:        218,
		Message:     "This is fine (Apache Web Server)",
		Description: "Used as a catch-all error condition for allowing response bodies to flow through Apache when ProxyErrorOverride is enabled.",
	},
	226: {
		Code:        226,
		Message:     "IM Used",
		Description: "The server has fulfilled a request for the resource, and the response is a representation of the result of one or more instance-manipulations applied to the current instance.",
	},
	300: {
		Code:        300,
		Message:     "Multiple Choices",
		Description: "A link list. The user can select a link and go to that location. Maximum five addresses.",
	},
	301: {
		Code:        301,
		Message:     "Moved Permanently",
		Description: "The requested page has moved to a new URL.",
	},
	302: {
		Code:        302,
		Message:     "Found",
		Description: "The requested page has moved temporarily to a new URL.",
	},
	303: {
		Code:        303,
		Message:     "See Other",
		Description: "The requested page can be found under a different URL.",
	},
	304: {
		Code:        304,
		Message:     "Not Modified",
		Description: "Indicates the requested page has not been modified since last requested.",
	},
	306: {
		Code:        306,
		Message:     "Switch Proxy",
		Description: "No longer used. Originally meant \"Subsequent requests should use the specified proxy.\"",
	},
	307: {
		Code:        307,
		Message:     "Temporary Redirect",
		Description: "The requested page has moved temporarily to a new URL.",
	},
	308: {
		Code:        308,
		Message:     "Resume Incomplete",
		Description: "Used in the resumable requests proposal to resume aborted PUT or POST requests.",
	},
	400: {
		Code:        400,
		Message:     "Bad Request",
		Description: "The request cannot be fulfilled due to bad syntax.",
	},
	401: {
		Code:        401,
		Message:     "Unauthorized",
		Description: "The request was a legal request, but the server is refusing to respond to it. For use when authentication is possible but has failed or not yet been provided.",
	},
	402: {
		Code:        402,
		Message:     "Payment Required",
		Description: "Not yet implemented by RFC standards, but reserved for future use.",
	},
	403: {
		Code:        403,
		Message:     "Forbidden",
		Description: "The request was a legal request, but the server is refusing to respond to it.",
	},
	404: {
		Code:        404,
		Message:     "Not Found",
		Description: "The requested page could not be found but may be available again in the future.",
	},
	405: {
		Code:        405,
		Message:     "Method Not Allowed",
		Description: "A request was made of a page using a request method not supported by that page.",
	},
	406: {
		Code:        406,
		Message:     "Not Acceptable",
		Description: "The server can only generate a response that is not accepted by the client.",
	},
	407: {
		Code:        407,
		Message:     "Proxy Authentication Required",
		Description: "The client must first authenticate itself with the proxy.",
	},
	408: {
		Code:        408,
		Message:     "Request Timeout",
		Description: "The server timed out waiting for the request.",
	},
	409: {
		Code:        409,
		Message:     "Conflict",
		Description: "The request could not be completed because of a conflict in the request.",
	},
	410: {
		Code:        410,
		Message:     "Gone",
		Description: "The requested page is no longer available.",
	},
	411: {
		Code:        411,
		Message:     "Length Required",
		Description: "The \"Content-Length\" is not defined. The server will not accept the request without it.",
	},
	412: {
		Code:        412,
		Message:     "Precondition Failed",
		Description: "The precondition given in the request evaluated to false by the server.",
	},
	413: {
		Code:        413,
		Message:     "Request Entity Too Large",
		Description: "The server will not accept the request, because the request entity is too large.",
	},
	414: {
		Code:        414,
		Message:     "Request-URI Too Long",
		Description: "The server will not accept the request, because the URL is too long. Occurs when you convert a POST request to a GET request with a long query information.",
	},
	415: {
		Code:        415,
		Message:     "Unsupported Media Type",
		Description: "The server will not accept the request, because the media type is not supported.",
	},
	416: {
		Code:        416,
		Message:     "Requested Range Not Satisfiable",
		Description: "The client has asked for a portion of the file, but the server cannot supply that portion.",
	},
	417: {
		Code:        417,
		Message:     "Expectation Failed",
		Description: "The server cannot meet the requirements of the Expect request-header field.",
	},
	418: {
		Code:        418,
		Message:     "I'm a teapot",
		Description: "Any attempt to brew coffee with a teapot should result in the error code \"418 I'm a teapot\". The resulting entity body MAY be short and stout.",
	},
	419: {
		Code:        419,
		Message:     "Page Expired (Laravel Framework)",
		Description: "Used by the Laravel Framework when a CSRF Token is missing or expired.",
	},
	420: {
		Code:        420,
		Message:     "Method Failure (Spring Framework)",
		Description: "A deprecated response used by the Spring Framework when a method has failed.",
	},
	421: {
		Code:        421,
		Message:     "Misdirected Request",
		Description: "The request was directed at a server that is not able to produce a response (for example because of connection reuse).",
	},
	422: {
		Code:        422,
		Message:     "Unprocessable Entity",
		Description: "The request was well-formed but was unable to be followed due to semantic errors.",
	},
	423: {
		Code:        423,
		Message:     "Locked",
		Description: "The resource that is being accessed is locked.",
	},
	424: {
		Code:        424,
		Message:     "Failed Dependency",
		Description: "The request failed due to failure of a previous request (e.g., a PROPPATCH).",
	},
	426: {
		Code:        426,
		Message:     "Upgrade Required",
		Description: "The client should switch to a different protocol.",
	},
	428: {
		Code:        428,
		Message:     "Precondition Required",
		Description: "The origin server requires the request to be conditional.",
	},
	429: {
		Code:        429,
		Message:     "Too Many Requests",
		Description: "The user has sent too many requests in a given amount of time.",
	},
	431: {
		Code:        431,
		Message:     "Request Header Fields Too Large",
		Description: "The server is unwilling to process the request because either an individual header field, or all the header fields collectively, are too large.",
	},
	440: {
		Code:        440,
		Message:     "Login Time-out (Microsoft)",
		Description: "The client's session has expired and must log in again.",
	},
	444: {
		Code:        444,
		Message:     "No Response (Nginx)",
		Description: "Used to indicate that the server has returned no information to the client and closed the connection.",
	},
	449: {
		Code:        449,
		Message:     "Retry With (Microsoft)",
		Description: "The request should be retried after performing the appropriate action.",
	},
	450: {
		Code:        450,
		Message:     "Blocked by Windows Parental Controls (Microsoft)",
		Description: "This error is given when Windows Parental Controls are turned on and are blocking access to the given webpage.",
	},
	451: {
		Code:        451,
		Message:     "Unavailable For Legal Reasons",
		Description: "The server is denying access to the resource as a consequence of a legal demand.",
	},
	460: {
		Code:        460,
		Message:     "Client closed connection (AWS ELB)",
		Description: "Used in AWS ELB when the client closes the connection.",
	},
	463: {
		Code:        463,
		Message:     "X-Forwarded-For header malformed (AWS ELB)",
		Description: "Used in AWS ELB when the X-Forwarded-For header is malformed.",
	},
	494: {
		Code:        494,
		Message:     "Request Header Too Large (Nginx)",
		Description: "Used in Nginx when the request header is too large.",
	},
	495: {
		Code:        495,
		Message:     "SSL Certificate Error (Nginx)",
		Description: "Used in Nginx when there is an SSL Certificate error.",
	},
	496: {
		Code:        496,
		Message:     "SSL Certificate Required (Nginx)",
		Description: "Used in Nginx when an SSL Certificate is required.",
	},
	497: {
		Code:        497,
		Message:     "HTTP Request Sent to HTTPS Port (Nginx)",
		Description: "Used in Nginx when an HTTP request is sent to an HTTPS port.",
	},
	498: {
		Code:        498,
		Message:     "Invalid Token (Esri)",
		Description: "Used in ArcGIS for Server when a token is invalid.",
	},
	499: {
		Code:        499,
		Message:     "Client Closed Request (Nginx)",
		Description: "Used in Nginx when the client closes the request.",
	},
	500: {
		Code:        500,
		Message:     "Internal Server Error",
		Description: "The request was not completed. The server met an unexpected condition.",
	},
	501: {
		Code:        501,
		Message:     "Not Implemented",
		Description: "The request was not completed. The server did not support the functionality required.",
	},
	502: {
		Code:        502,
		Message:     "Bad Gateway",
		Description: "The request was not completed. The server received an invalid response from the upstream server.",
	},
	503: {
		Code:        503,
		Message:     "Service Unavailable",
		Description: "The request was not completed. The server is temporarily overloading or down.",
	},
	504: {
		Code:        504,
		Message:     "Gateway Timeout",
		Description: "The gateway has timed out.",
	},
	505: {
		Code:        505,
		Message:     "HTTP Version Not Supported",
		Description: "The server does not support the \"HTTP protocol\" version.",
	},
	506: {
		Code:        506,
		Message:     "Variant Also Negotiates",
		Description: "The server has an internal configuration error: transparent content negotiation for the request results in a circular reference.",
	},
	507: {
		Code:        507,
		Message:     "Insufficient Storage",
		Description: "The server is unable to store the representation needed to complete the request.",
	},
	508: {
		Code:        508,
		Message:     "Loop Detected",
		Description: "The server detected an infinite loop while processing the request.",
	},
	509: {
		Code:        509,
		Message:     "Bandwidth Limit Exceeded (Apache Web Server/cPanel)",
		Description: "Used by the Apache Web Server/cPanel when a bandwidth limit is exceeded.",
	},
	510: {
		Code:        510,
		Message:     "Not Extended",
		Description: "Further extensions to the request are required for the server to fulfil it.",
	},
	511: {
		Code:        511,
		Message:     "Network Authentication Required",
		Description: "The client needs to authenticate to gain network access.",
	},
	520: {
		Code:        520,
		Message:     "Unknown Error (Cloudflare)",
		Description: "Used as a catch-all response for Cloudflare when the origin server returns something unexpected.",
	},
	521: {
		Code:        521,
		Message:     "Web Server Is Down (Cloudflare)",
		Description: "Used when the origin server refuses connections from Cloudflare.",
	},
	522: {
		Code:        522,
		Message:     "Connection Timed Out (Cloudflare)",
		Description: "Used when Cloudflare is unable to reach the origin server in a timely manner.",
	},
	523: {
		Code:        523,
		Message:     "Origin Is Unreachable (Cloudflare)",
		Description: "Used when Cloudflare cannot reach the origin server.",
	},
	524: {
		Code:        524,
		Message:     "A Timeout Occurred (Cloudflare)",
		Description: "Used when Cloudflare times out contacting the origin server.",
	},
	525: {
		Code:        525,
		Message:     "SSL Handshake Failed (Cloudflare)",
		Description: "Used when Cloudflare fails to negotiate a TLS/SSL handshake with the origin server.",
	},
	526: {
		Code:        526,
		Message:     "Invalid SSL Certificate (Cloudflare)",
		Description: "Used when Cloudflare cannot validate the SSL certificate on the origin server.",
	},
	527: {
		Code:        527,
		Message:     "Railgun Error (Cloudflare)",
		Description: "Used when there is a Railgun error in the Cloudflare network.",
	},
	530: {
		Code:        530,
		Message:     "Site is Frozen",
		Description: "Used when the site is frozen due to inactivity or billing issues.",
	},
	598: {
		Code:        598,
		Message:     "Network Read Timeout Error",
		Description: "Used to indicate a network read timeout behind the proxy.",
	},
	599: {
		Code:        599,
		Message:     "Network Connect Timeout Error",
		Description: "Used to indicate a network connect timeout behind the proxy.",
	},
}

// For the given status code, returns a concatenated string of the code and the message i.e 200 OK
// This is used in constructing response string.
func GetStatusString(code int) string {

	status := HTTPStatuses[code]
	statusString := strconv.Itoa(status.Code) + " " + status.Message
	return statusString
}

// Descriptive names for HTTP Status codes.
const (
	HTTP_100_CONTINUE                      = 100
	HTTP_101_SWITCHING_PROTOCOLS           = 101
	HTTP_102_PROCESSING                    = 102
	HTTP_103_EARLY_HINTS                   = 103
	HTTP_200_OK                            = 200
	HTTP_201_CREATED                       = 201
	HTTP_202_ACCEPTED                      = 202
	HTTP_203_NON_AUTHORITATIVE_INFORMATION = 203
	HTTP_204_NO_CONTENT                    = 204
	HTTP_205_RESET_CONTENT                 = 205
	HTTP_206_PARTIAL_CONTENT               = 206
	HTTP_207_MULTI_STATUS                  = 207
	HTTP_208_ALREADY_REPORTED              = 208
	HTTP_226_IM_USED                       = 226

	HTTP_300_MULTIPLE_CHOICES   = 300
	HTTP_301_MOVED_PERMANENTLY  = 301
	HTTP_302_FOUND              = 302
	HTTP_303_SEE_OTHER          = 303
	HTTP_304_NOT_MODIFIED       = 304
	HTTP_305_USE_PROXY          = 305
	HTTP_306_RESERVED           = 306
	HTTP_307_TEMPORARY_REDIRECT = 307
	HTTP_308_PERMANENT_REDIRECT = 308
	HTTP_400_BAD_REQUEST        = 400

	HTTP_401_UNAUTHORIZED                    = 401
	HTTP_402_PAYMENT_REQUIRED                = 402
	HTTP_403_FORBIDDEN                       = 403
	HTTP_404_NOT_FOUND                       = 404
	HTTP_405_METHOD_NOT_ALLOWED              = 405
	HTTP_406_NOT_ACCEPTABLE                  = 406
	HTTP_407_PROXY_AUTHENTICATION_REQUIRED   = 407
	HTTP_408_REQUEST_TIMEOUT                 = 408
	HTTP_409_CONFLICT                        = 409
	HTTP_410_GONE                            = 410
	HTTP_411_LENGTH_REQUIRED                 = 411
	HTTP_412_PRECONDITION_FAILED             = 412
	HTTP_413_REQUEST_ENTITY_TOO_LARGE        = 413
	HTTP_414_REQUEST_URI_TOO_LONG            = 414
	HTTP_415_UNSUPPORTED_MEDIA_TYPE          = 415
	HTTP_416_REQUESTED_RANGE_NOT_SATISFIABLE = 416
	HTTP_417_EXPECTATION_FAILED              = 417
	HTTP_418_IM_A_TEAPOT                     = 418
	HTTP_421_MISDIRECTED_REQUEST             = 421
	HTTP_422_UNPROCESSABLE_ENTITY            = 422
	HTTP_423_LOCKED                          = 423
	HTTP_424_FAILED_DEPENDENCY               = 424
	HTTP_425_TOO_EARLY                       = 425
	HTTP_426_UPGRADE_REQUIRED                = 426
	HTTP_428_PRECONDITION_REQUIRED           = 428
	HTTP_429_TOO_MANY_REQUESTS               = 429
	HTTP_431_REQUEST_HEADER_FIELDS_TOO_LARGE = 431
	HTTP_451_UNAVAILABLE_FOR_LEGAL_REASONS   = 451

	HTTP_500_INTERNAL_SERVER_ERROR           = 500
	HTTP_501_NOT_IMPLEMENTED                 = 501
	HTTP_502_BAD_GATEWAY                     = 502
	HTTP_503_SERVICE_UNAVAILABLE             = 503
	HTTP_504_GATEWAY_TIMEOUT                 = 504
	HTTP_505_HTTP_VERSION_NOT_SUPPORTED      = 505
	HTTP_506_VARIANT_ALSO_NEGOTIATES         = 506
	HTTP_507_INSUFFICIENT_STORAGE            = 507
	HTTP_508_LOOP_DETECTED                   = 508
	HTTP_509_BANDWIDTH_LIMIT_EXCEEDED        = 509
	HTTP_510_NOT_EXTENDED                    = 510
	HTTP_511_NETWORK_AUTHENTICATION_REQUIRED = 511
)
