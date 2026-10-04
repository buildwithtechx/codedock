export function printSuccess(message: string): void {
  console.log(`\x1b[32m✔\x1b[0m ${message}`);
}

export function printError(message: string): void {
  console.error(`\x1b[31m✖\x1b[0m ${message}`);
}

export function printInfo(message: string): void {
  console.log(`\x1b[36mℹ\x1b[0m ${message}`);
}

export function printJson(data: unknown): void {
  console.log(JSON.stringify(data, null, 2));
}

export function printTable(
  headers: string[],
  rows: (string | number | boolean | undefined)[][]
): void {
  const colWidths = headers.map((header, colIndex) => {
    let maxWidth = header.length;
    for (const row of rows) {
      const val = row[colIndex] === undefined ? '' : String(row[colIndex]);
      if (val.length > maxWidth) {
        maxWidth = val.length;
      }
    }
    return maxWidth;
  });

  const headerRow = headers.map((header, index) => header.padEnd(colWidths[index])).join('  ');
  console.log(`\x1b[1m${headerRow}\x1b[0m`);

  const separator = colWidths.map((width) => '─'.repeat(width)).join('  ');
  console.log(`\x1b[90m${separator}\x1b[0m`);

  for (const row of rows) {
    const formatted = row
      .map((cell, index) => {
        const val = cell === undefined ? '' : String(cell);
        return val.padEnd(colWidths[index]);
      })
      .join('  ');
    console.log(formatted);
  }
}
