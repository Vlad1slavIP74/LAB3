const fs = require('fs');


for (let i = 0; i < 1000000; i++) {
    fs.writeFile(`./test/mynewfile${i}.txt`, 'Hello content!', function (err) {
        if (err) throw err;
      });
}

