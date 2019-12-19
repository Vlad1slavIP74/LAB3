'use strict'

const fs = require('fs')
const path = require('path');

function countFileLines (fileName, sourceDir, destDir) {
    const filePath = `${sourceDir}/${fileName}`
    return new Promise((resolve, reject) => {
    let lineCount = 0;
    fs.createReadStream(filePath)
      .on("data", (buffer) => {
        let idx = -1;
        do {
          idx = buffer.indexOf(10, idx+1);
          lineCount++;
        } while (idx !== -1);
      }).on("end", () => {
        const extension = path.extname(fileName);
        const newFileName = path.basename(fileName,extension);
        const pathToOutputFile = `${destDir}/${newFileName}.res`;
        fs.writeFile(pathToOutputFile, lineCount, (err) => {
          if (err) reject(err);
          resolve();
        });
        resolve(lineCount);
      }).on("error", reject);
    });
  };


(async function main(){
    let pathToSourceDir;
    let pathToDestinationDir;
    try {
        pathToSourceDir = path.resolve(process.argv[2]);
      } catch (err) {
        console.error('Wrong source directory path!');
        process.exit(1);
      }
    
      try {
        pathToDestinationDir = path.resolve(process.argv[3]);
      } catch (err) {
        console.error('Wrong destination directory path!');
      }
      if(!fs.existsSync(path.resolve(process.argv[3]))) {
          console.warn('Create dir\t' +process.argv[3])
          fs.mkdirSync(path.resolve(process.argv[3]))
        }
         
    const inputFiles = fs.readdirSync(pathToSourceDir);
    const promise = [];
    for (const file of inputFiles) {
        const res = countFileLines(file, pathToSourceDir, pathToDestinationDir);
        promise.push(res)
    }
    await Promise.all(promise);
})();