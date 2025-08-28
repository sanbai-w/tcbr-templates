const path = require("path");
const express = require("express");
const cors = require("cors");
const morgan = require("morgan");
const { init: initDB, Counter, checkConnection } = require("./db");

const logger = morgan("tiny");

const app = express();
app.use(express.urlencoded({ extended: false }));
app.use(express.json());
app.use(cors());
app.use(logger);

// 首页
app.get("/", async (req, res) => {
  res.sendFile(path.join(__dirname, "index.html"));
});

// 获取数据库状态
app.get("/api/status", async (req, res) => {
  res.send({
    code: 0,
    data: {
      connected: checkConnection(),
    },
  });
});

// 更新计数
app.post("/api/count", async (req, res) => {
  if (!checkConnection()) {
    return res.status(503).send({
      code: -1,
      message: "数据库未连接，请检查配置",
      data: null,
    });
  }

  try {
    const { action } = req.body;
    if (action === "inc") {
      await Counter.create();
    } else if (action === "clear") {
      await Counter.destroy({
        truncate: true,
      });
    }
    const result = await Counter.count();
    res.send({
      code: 0,
      data: result,
    });
  } catch (error) {
    res.status(500).send({
      code: -1,
      message: "数据库操作失败",
      data: null,
    });
  }
});

// 获取计数
app.get("/api/count", async (req, res) => {
  if (!checkConnection()) {
    return res.status(503).send({
      code: -1,
      message: "数据库未连接，请检查配置",
      data: null,
    });
  }

  try {
    const result = await Counter.count();
    res.send({
      code: 0,
      data: result,
    });
  } catch (error) {
    res.status(500).send({
      code: -1,
      message: "数据库操作失败",
      data: null,
    });
  }
});

// 小程序调用，获取微信 Open ID
app.get("/api/wx_openid", async (req, res) => {
  if (req.headers["x-wx-source"]) {
    res.send(req.headers["x-wx-openid"]);
  }
});

const port = process.env.PORT || 8080;

async function bootstrap() {
  try {
    await initDB();
    app.listen(port, () => {
      console.log("启动成功", port);
      if (!checkConnection()) {
        console.log("⚠️  警告: 数据库连接失败，应用将以有限功能模式运行");
        console.log("请检查 .env 文件中的数据库配置");
      }
    });
  } catch (error) {
    console.error("启动失败:", error);
    process.exit(1);
  }
}

bootstrap();
