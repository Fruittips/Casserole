/**
  A centralised file for all routes for the handlers

NOTE: If this is changed, routes in HttpClient.go need to be changed too!
*/
package handlers

//TODO: Check if write endpoints actually need the student parameter

const READ_ENDPOINT_FSTRING = "/read/course/%v/student/%v"
const WRITE_ENDPOINT_FSTRING = "/write/course/%v/student/%v"
const INTERNAL_READ_ENDPOINT_FSTRING = "/internal/read/course/%v/student/%v"
const INTERNAL_WRITE_ENDPOINT_FSTRING = "/internal/write/course/%v/student/%v"
const INTERNAL_KILL_ENDPOINT_FSTRING = "/internal/kill"
const INTERNAL_REVIVE_ENDPOINT_FSTRING = "/internal/revive"
const INTERNAL_CHECKHH_ENDPOINT_FSTRING = "/internal/checkhh/node/%v"
