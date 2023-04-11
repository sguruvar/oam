const express = require('express');
const bodyParser = require('body-parser');
const cors = require('cors');
const axios = require('axios');

const app = express();
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: true }));
app.use(cors());
app.use(express.static('public'))

APP_URL='http://backend:8080'

// Endpoint for submitting a new guestbook entry
app.post('/submit', (req, res) => {
    const name = req.body.name;
    const message = req.body.message;

    // Make a POST request to the /submit endpoint
    axios.post(`${APP_URL}/submitentry`, { name, message })
        .then(response => {
            res.sendStatus(200);
        })
        .catch(error => {
            console.log(error);
            res.sendStatus(500);
        });
});

// Endpoint for getting all guestbook entries
app.get('/getEntries', (req, res) => {

    // Make a GET request to the /getentries endpoint
    axios.get(`${APP_URL}/getentries`)
        .then(response => {
            res.send(response.data);
        })
        .catch(error => {
            console.log(error);
            res.sendStatus(500);
        });
});

// Start the server
const port = 3000;
app.listen(port, () => {
    console.log(`Server running on port ${port}`);
});