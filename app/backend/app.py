from flask import Flask, render_template, request, json,jsonify

app = Flask(__name__)

entries = []

@app.route('/')
def index():
    return 'hello from backend'

@app.route('/submitentry', methods=['POST'])
def submit_entry():
    name = request.json['name']
    message = request.json['message']
    entries.append({'name': name, 'message': message})
    return jsonify({'entries': entries})

@app.route('/getentries')
def getentries():
    response = app.response_class(
        response=json.dumps(entries),
        status=200,
        mimetype='application/json'
    )
    return response


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8080, debug=True)