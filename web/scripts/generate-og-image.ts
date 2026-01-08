import { chromium } from '@playwright/test';
import { writeFileSync } from 'fs';
import { join } from 'path';

const html = `
<!DOCTYPE html>
<html>
<head>
  <style>
    * { margin: 0; padding: 0; box-sizing: border-box; }
    body {
      width: 1200px;
      height: 630px;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      background: linear-gradient(135deg, #0078d4 0%, #004578 100%);
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', system-ui, sans-serif;
      color: white;
      padding: 60px;
    }
    .container {
      text-align: center;
      max-width: 1000px;
    }
    .title {
      font-size: 96px;
      font-weight: 700;
      margin-bottom: 30px;
      letter-spacing: -2px;
    }
    .command {
      font-size: 48px;
      font-family: 'Consolas', 'Monaco', monospace;
      background: rgba(0, 0, 0, 0.3);
      padding: 20px 40px;
      border-radius: 12px;
      margin-bottom: 40px;
      display: inline-block;
    }
    .description {
      font-size: 36px;
      opacity: 0.95;
      line-height: 1.4;
      font-weight: 400;
    }
    .logo {
      position: absolute;
      bottom: 40px;
      right: 40px;
      font-size: 24px;
      opacity: 0.8;
    }
  </style>
</head>
<body>
  <div class="container">
    <div class="title">azd exec</div>
    <div class="command">azd exec ./script.sh</div>
    <div class="description">Execute scripts with your azd env context!</div>
  </div>
  <div class="logo">jongio.github.io/azd-exec</div>
</body>
</html>
`;

async function generateOgImage() {
  const browser = await chromium.launch();
  const page = await browser.newPage({
    viewport: { width: 1200, height: 630 }
  });
  
  await page.setContent(html);
  const screenshot = await page.screenshot({ type: 'png' });
  
  const outputPath = join(process.cwd(), 'public', 'og-image.png');
  writeFileSync(outputPath, screenshot);
  
  console.log('✓ OG image generated:', outputPath);
  
  await browser.close();
}

generateOgImage().catch(console.error);
