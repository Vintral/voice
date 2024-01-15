import http from 'k6/http';
import { check, sleep } from "k6";

export let options = {
  thresholds: {
    http_req_failed: ['rate<0.01'], // http errors should be less than 1%
    http_req_duration: ['p(95)<50'], // 95% of requests should be below 50ms
  },
  vus: 250,
  duration: "1h",
};

const priorities = [
  { priority: 200, type: "public_issues" },
  { priority: 200, type: "admin_issues" },
  { priority: 200, type: "health" },
  { priority: 200, type: "get_user" },
  { priority: 200, type: "update_user" },
  { priority: 200, type: "get_donations" },
  { priority: 200, type: "create_donation" },
  { priority: 200, type: "delete_donation" },
]

var defaultParams = {};
var submitParams = {}
var issue = "";
var token = "";
var added = [];

export function setup() {
  token = retrieveToken();  
  issue = retrieveIssue();
  
  return { token, issue }
}

export default function ( data ) {
  issue = data.issue;
  token = data.token;
  defaultParams = {
    headers: {
      "Authorization": "Bearer " + token
    }
  }

  submitParams = {
    headers: {
      "Content-Type": "application/json",
      "Authorization": "Bearer " + token
    }
  }
  
  const totalPriority = priorities.reduce( ( prev, curr ) => {
    return prev + curr.priority;
  }, 0 );
  let picker = Math.floor( Math.random() * totalPriority );

  let req;
  priorities
    .sort( sorter )
    .every( request => {      
      picker -= request.priority;
      if( picker < 0 ) {        
        req = request;
        return false;
      } else return true;
    } );

  switch( req.type ) {
    case "public_issues": getPublicIssues(); break;
    case "admin_issues": getAdminIssues(); break;
    case "health": sendHealthCheck(); break;
    case "get_user": getUser(); break;
    case "update_user": updateUser(); break;
    case "get_donations": getDonations(); break;
    case "create_donation": createDonation(); break;
    case "delete_donation": deleteDonation(); break;
  }

  sleep( ( Math.floor( Math.random() * 1000 ) + 100 ) / 1000 );
};

const sorter = ( a, b ) => {
  if( a.priority < b.priority ) return -1;
  if( a.priority > b.priority ) return 1;
  return 0;
}

const retrieveToken = () => {
  const res = http.get( `${__ENV.API_URL}/auth`)
  return res.body;
}

const retrieveIssue = () => {
  const res = http.get( `${__ENV.API_URL}/public/issues` );
  const data = JSON.parse( res.body );

  return data[ 0 ].guid;
}

const sendHealthCheck = () => {
  const response = http.get( `${__ENV.API_URL}` );
  check(response, { "(GET)/": (r) => r.status === 200 });  
}

const getPublicIssues = () => {
  const response = http.get( `${__ENV.API_URL}/public/issues` );
  check(response, { "(GET)/public/issues": (r) => r.status === 200 });  
}

const getAdminIssues = () => {    
  const response = http.get( `${__ENV.API_URL}/admin/issues`, defaultParams );
  check(response, { "(GET)/admin/issues": (r) => r.status === 200 });  
}

const getUser = () => {  
  const response = http.get( `${__ENV.API_URL}/user`, defaultParams );
  check(response, { "(GET)/user": (r) => r.status === 200 });  
}

const updateUser = () => {
  const body = {
    "first_name": "Jane",
    "last_name": "Doe",
    "address_1": "123 Nowhere St",
    "address_2": "Apt C",
    "city": "Somewhere",
    "state": "MO",
    "zip": "12345"
  }

  const response = http.patch( `${__ENV.API_URL}/user`, JSON.stringify( body ), submitParams );
  check(response, { "(PATCH)/user": (r) => r.status === 200 });  
}

const getDonations = () => {
  const response = http.get( `${__ENV.API_URL}/user/donations`, defaultParams );
  check(response, { "(GET)/user/donations": (r) => r.status === 200 });  
}

const createDonation = () => {
  const body = {
    issue,
    amount: ( Math.floor( Math.random() * 1000 ) ) / 100,
  }

  const response = http.post( `${__ENV.API_URL}/user/donation`, JSON.stringify( body ), submitParams );
  check(response, { "(POST)/user/donation": (r) => r.status === 201 });  
  
  const ret = JSON.parse( response.body );
  added.push( ret.guid );
}

const deleteDonation = () => {
  const donation = added.shift();
  if( !donation ) return;  

  const response = http.del( `${__ENV.API_URL}/user/donation/${donation}`, null, defaultParams );
  check(response, { "(DELETE)/user/donation": (r) => r.status === 200 });    
}