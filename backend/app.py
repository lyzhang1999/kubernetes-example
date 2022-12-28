from flask import Flask
from flask_sqlalchemy import SQLAlchemy
import os
import subprocess
from flask import request

app = Flask(__name__)

# This keeps the warning messages away
app.config['SQLALCHEMY_TRACK_MODIFICATIONS'] = False
db_uri = os.environ.get('DATABASE_URI')
db_username = os.environ.get('DATABASE_USERNAME')
db_password = os.environ.get('DATABASE_PASSWORD')
app.config['SQLALCHEMY_DATABASE_URI'] = "postgresql://"+ db_username +":"+ db_password +"@" + \
    db_uri+"/postgres"

db = SQLAlchemy(app)
db.init_app(app)


class Text(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    text = db.Column(db.String(120), nullable=False)

    def __init__(self, text):
        self.text = text

    def __repr__(self):
        return f"<Text {self.text}>"


@app.route('/healthy')
def healthy():
    return {"healthy": True}, 200


@app.route('/host_name')
def host_name():
    return {"host_name": os.environ.get('HOSTNAME')}, 200


@app.route('/ab')
def ab():
    subprocess.run(["ab", "-c", "50", "-n", "10000", "http://127.0.0.1:5000/"])
    return {"ab": True}, 200


@app.route('/fetch')
def fetch():
    words = Text.query.all()
    results = [
        {
            "text": word.text
        } for word in words]
    return {"texts": results}, 200


@app.route('/add', methods=['POST'])
def add():
    text = request.json['text']
    db.session.add(Text(text=text))
    db.session.commit()
    return 'Done', 201


@app.route('/delete', methods=['DELETE'])
def delete():
    db.session.query(Text).delete()
    db.session.commit()
    return 'Done', 200


if __name__ == "__main__":
    app.run(host='0.0.0.0')
