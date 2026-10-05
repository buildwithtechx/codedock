import { createInterface } from 'node:readline/promises';
import { Writable } from 'node:stream';

export async function readPassword(): Promise<string> {
  if (!process.stdin.isTTY)
    throw new Error(
      'Interactive login requires a terminal. Use an API token for noninteractive login.'
    );
  process.stdout.write('Password: ');
  const output = new Writable({
    write(_chunk, _encoding, callback) {
      callback();
    },
  });
  const prompt = createInterface({ input: process.stdin, output, terminal: true });
  try {
    return await new Promise<string>((resolve, reject) => {
      prompt.once('SIGINT', () => reject(new Error('Login cancelled')));
      prompt.question('').then(resolve, reject);
    });
  } finally {
    prompt.close();
    output.end();
    process.stdout.write('\n');
  }
}
