import path from "node:path";
import fs from "node:fs/promises";
import { simpleParser } from "mailparser";

const emailDir = "./emails";
const desiredDomain = "fotball.no";

async function main() {
  const files = await fs.readdir(emailDir, { recursive: true });

  for (const file of files) {
    let filePath = path.join(emailDir, file);
    if (filePath.endsWith(".pdf") || file.includes("old")) {
      continue;
    }
    const data = await Bun.file(filePath).text();

    try {
      const parsed = await simpleParser(data);
      const sender = parsed.from.value[0].address;
      // const recipient = parsed.to.value[0].address;

      const senderDomain = sender.split("@")[1].toLowerCase();

      if (senderDomain === desiredDomain) {
        console.log(`Email from ${sender} found in ${file}`);
      }

      // if (recipient === "kurlandfk@ebilag.catacloud.com") {
      //   console.log(parsed.text);
      // }
    } catch (parseErr) {
      console.error("Error parsing email:", parseErr);
    }
  }
}

main().catch(console.error);
