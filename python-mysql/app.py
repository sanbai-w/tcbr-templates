import os
from typing import Any, Dict

from flask import Flask, jsonify, request, send_file
from flask_cors import CORS
from sqlalchemy.exc import SQLAlchemyError

from db import Counter, SessionLocal, init_db, truncate_or_delete_all, check_connection


app = Flask(__name__, static_folder=None)
CORS(app)


@app.route("/", methods=["GET"])  # Home page
def index():
    html_path = os.path.join(os.path.dirname(__file__), "index.html")
    return send_file(html_path)


@app.get("/api/status")
def get_status():
    """Get database connection status."""
    return jsonify({"code": 0, "data": {"connected": check_connection()}})


@app.get("/api/count")
def get_count():
    if not check_connection():
        return (
            jsonify({"code": -1, "message": "数据库未连接，请检查配置", "data": None}),
            503,
        )

    session = SessionLocal()
    try:
        total_count: int = session.query(Counter).count()
        return jsonify({"code": 0, "data": total_count})
    except SQLAlchemyError as e:
        return jsonify({"code": -1, "message": "数据库操作失败", "data": None}), 500
    finally:
        session.close()


@app.post("/api/count")
def modify_count():
    if not check_connection():
        return (
            jsonify({"code": -1, "message": "数据库未连接，请检查配置", "data": None}),
            503,
        )

    session = SessionLocal()
    try:
        payload: Dict[str, Any] = request.get_json(silent=True) or {}
        action = payload.get("action")

        if action == "inc":
            session.add(Counter(count=1))
            session.commit()
        elif action == "clear":
            truncate_or_delete_all(session)

        total_count: int = session.query(Counter).count()
        return jsonify({"code": 0, "data": total_count})
    except SQLAlchemyError as e:
        session.rollback()
        return jsonify({"code": -1, "message": "数据库操作失败", "data": None}), 500
    finally:
        session.close()


@app.get("/api/wx_openid")
def get_wx_openid():
    # Mirror Node.js behavior: only respond when header x-wx-source exists
    if request.headers.get("x-wx-source"):
        return request.headers.get("x-wx-openid", "")
    return ("", 204)


def bootstrap() -> Flask:
    try:
        init_db()
        return app
    except Exception as e:
        print(f"启动失败: {str(e)}")
        # Still return app even if database initialization fails
        return app


app = bootstrap()


if __name__ == "__main__":
    port = int(os.getenv("PORT", "8080"))
    print(f"启动成功 {port}")
    if not check_connection():
        print("⚠️  警告: 数据库连接失败，应用将以有限功能模式运行")
        print("请检查 .env 文件中的数据库配置")
    app.run(host="0.0.0.0", port=port)
