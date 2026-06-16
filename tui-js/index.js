import { select, input } from '@inquirer/prompts';
import crypto from 'node:crypto';

function clearScreen() {
  console.clear();
}

function showHeader() {
  console.log('╔════════════════════════════╗');
  console.log('║        MY TERMINAL         ║');
  console.log('╚════════════════════════════╝');
  console.log('');
}

function randomPassword(length) {
  const chars =
    'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*';

  let result = '';

  for (let i = 0; i < length; i++) {
    const index = crypto.randomInt(0, chars.length);
    result += chars[index];
  }

  return result;
}

async function waitEnter() {
  await input({
    message: '按回车返回首页',
  });
}

async function passwordPage() {
  clearScreen();
  showHeader();

  const lengthText = await input({
    message: '请输入密码长度：',
    default: '16',
  });

  const length = Number(lengthText);

  if (!Number.isInteger(length) || length <= 0) {
    console.log('');
    console.log('密码长度必须是正整数。');
    await waitEnter();
    return;
  }

  const password = randomPassword(length);

  console.log('');
  console.log('生成的密码：');
  console.log(password);
  console.log('');

  await waitEnter();
}

async function aboutPage() {
  clearScreen();
  showHeader();

  console.log('这是一个最小版 npm 终端交互面板。');
  console.log('目前功能：');
  console.log('- 首页菜单');
  console.log('- 密码生成器');
  console.log('- 关于页面');
  console.log('- 退出');
  console.log('');

  await waitEnter();
}

async function main() {
  while (true) {
    clearScreen();
    showHeader();

    const choice = await select({
      message: '请选择功能：',
      choices: [
        {
          name: '密码生成器',
          value: 'password',
        },
        {
          name: '关于',
          value: 'about',
        },
        {
          name: '退出',
          value: 'exit',
        },
      ],
    });

    if (choice === 'password') {
      await passwordPage();
    }

    if (choice === 'about') {
      await aboutPage();
    }

    if (choice === 'exit') {
      clearScreen();
      console.log('再见！');
      process.exit(0);
    }
  }
}

main();