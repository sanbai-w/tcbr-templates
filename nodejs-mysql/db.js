// 加载 .env 文件中的环境变量
require("dotenv").config();

const { Sequelize, DataTypes } = require("sequelize");

// 从环境变量中读取数据库配置
const {
  MYSQL_USERNAME,
  MYSQL_PASSWORD,
  MYSQL_ADDRESS = "",
  MYSQL_DATABASE = "",
  TABLE_NAME = "counters",
} = process.env;

const [host, port] = MYSQL_ADDRESS.split(":");

const sequelize = new Sequelize(
  MYSQL_DATABASE,
  MYSQL_USERNAME,
  MYSQL_PASSWORD,
  {
    host,
    port,
    dialect: "mysql" /* one of 'mysql' | 'mariadb' | 'postgresql' | 'mssql' */,
  }
);

// 定义数据模型
const Counter = sequelize.define(TABLE_NAME, {
  count: {
    type: DataTypes.INTEGER,
    allowNull: false,
    defaultValue: 1,
  },
});

// 数据库连接状态
let isConnected = false;

// 数据库初始化方法
async function init() {
  try {
    // 测试数据库连接
    await sequelize.authenticate();
    await Counter.sync({ alter: true });
    isConnected = true;
    console.log("数据库连接成功");
  } catch (error) {
    isConnected = false;
    console.error("数据库连接失败:", error.message);
    console.log("请检查 .env 文件中的数据库配置");
  }
}

// 检查数据库连接状态
function checkConnection() {
  return isConnected;
}

// 导出初始化方法和模型
module.exports = {
  init,
  Counter,
  checkConnection,
};
