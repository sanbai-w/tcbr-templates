import os
from typing import Tuple
from dotenv import load_dotenv
from datetime import datetime

# Load environment variables from .env file
load_dotenv()

from sqlalchemy import Column, Integer, create_engine, text, DateTime
from sqlalchemy.orm import declarative_base, sessionmaker
from sqlalchemy.exc import SQLAlchemyError

# Read database configuration from environment variables
MYSQL_USERNAME = os.getenv("MYSQL_USERNAME", "")
MYSQL_PASSWORD = os.getenv("MYSQL_PASSWORD", "")
MYSQL_ADDRESS = os.getenv("MYSQL_ADDRESS", "")
MYSQL_DATABASE = os.getenv("MYSQL_DATABASE", "")
TABLE_NAME = os.getenv("TABLE_NAME", "Counter")

# Parse host and port from MYSQL_ADDRESS
# Fallback to default MySQL port 3306 if not provided
host: str = ""
port: str = "3306"
if MYSQL_ADDRESS:
    parts: Tuple[str, ...] = tuple(MYSQL_ADDRESS.split(":"))
    if len(parts) >= 1:
        host = parts[0]
    if len(parts) >= 2 and parts[1]:
        port = parts[1]

# Create SQLAlchemy engine
DATABASE_URL = f"mysql+pymysql://{MYSQL_USERNAME}:{MYSQL_PASSWORD}@{host}:{port}/{MYSQL_DATABASE}?charset=utf8mb4"
engine = create_engine(
    DATABASE_URL,
    pool_pre_ping=True,  # avoid stale connections
    pool_recycle=3600,
    future=True,
)

# Base class for models
Base = declarative_base()

# Database connection status
is_connected = False


def create_counter_model(table_name: str):
    """Create a dynamic Counter model with a runtime table name.

    The table will contain an auto-increment primary key and a single
    integer column named `count` defaulting to 1, mirroring the Node.js demo.
    """

    class Counter(Base):
        __tablename__ = table_name

        id = Column(Integer, primary_key=True, autoincrement=True)
        count = Column(Integer, nullable=False, default=1)
        createdAt = Column(DateTime, nullable=False, default=datetime.now)
        updatedAt = Column(DateTime, nullable=False, default=datetime.now)

    Counter.__name__ = f"{table_name}Model"
    return Counter


# Expose the dynamic model using the configured table name
Counter = create_counter_model(TABLE_NAME)

# Session factory
SessionLocal = sessionmaker(bind=engine, autocommit=False, autoflush=False, future=True)


def init_db() -> None:
    """Initialize database schema if not exists."""
    global is_connected
    try:
        # Test database connection
        with engine.connect() as connection:
            connection.execute(text("SELECT 1"))

        # Create tables if connection successful
        Base.metadata.create_all(engine)
        is_connected = True
        print("数据库连接成功")
    except SQLAlchemyError as e:
        is_connected = False
        print(f"数据库连接失败: {str(e)}")
        print("请检查 .env 文件中的数据库配置")


def check_connection() -> bool:
    """Check database connection status."""
    return is_connected


def truncate_or_delete_all(session) -> None:
    """Clear all rows in the Counter table using TRUNCATE first, falling back to DELETE.

    This mirrors the Node.js demo's `truncate: true` behavior when clearing.
    """
    try:
        session.execute(text(f"TRUNCATE TABLE `{TABLE_NAME}`"))
        session.commit()
    except SQLAlchemyError:
        session.rollback()
        session.query(Counter).delete()
        session.commit()
