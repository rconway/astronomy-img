import { promises as fs } from "fs";
// import fs_std from "fs";
// const fs = fs_std.promises;
import fetch from "node-fetch";

const IMG_FILE = "astronomy-img.jpg"

async function main() {
    const response = await fetch("https://api.nasa.gov/planetary/apod?api_key=DEMO_KEY")
    const json = await response.json()
    const imgResponse = await fetch(json.url)
    const img = await imgResponse.arrayBuffer()
    await fs.writeFile(IMG_FILE, Buffer.from(img))
    console.log(`\nImage Description: ${json.explanation}`)
    console.log(`\nImage written to file: ${IMG_FILE}\n`)
}

main()
