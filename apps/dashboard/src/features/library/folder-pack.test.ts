import { describe, expect, it } from 'vitest';
import { type FolderEntry, inspectFolderFiles } from './folder-pack';

function makeFile(name: string, content = ''): File {
  return new File([content], name, { type: 'text/plain' });
}

describe('inspectFolderFiles', () => {
  it('detects Static Site when index.html is present without backend manifests', async () => {
    const entries: FolderEntry[] = [
      { path: 'index.html', file: makeFile('index.html', '<h1>Hello</h1>') },
      { path: 'style.css', file: makeFile('style.css', 'body { color: red; }') },
      { path: 'app.js', file: makeFile('app.js', 'console.log("hi");') },
    ];

    const result = await inspectFolderFiles(entries);
    expect(result.framework).toBe('Static Site');
    expect(result.fileCount).toBe(3);
  });

  it('detects Go stack when go.mod is present', async () => {
    const entries: FolderEntry[] = [
      { path: 'go.mod', file: makeFile('go.mod', 'module myapp\ngo 1.22\n') },
      { path: 'main.go', file: makeFile('main.go', 'package main') },
    ];

    const result = await inspectFolderFiles(entries);
    expect(result.framework).toBe('Go');
  });

  it('detects Python stack when requirements.txt is present', async () => {
    const entries: FolderEntry[] = [
      { path: 'requirements.txt', file: makeFile('requirements.txt', 'flask\n') },
      { path: 'app.py', file: makeFile('app.py', 'print("hello")') },
    ];

    const result = await inspectFolderFiles(entries);
    expect(result.framework).toBe('Python');
  });

  it('detects Node stack from package.json', async () => {
    const entries: FolderEntry[] = [
      {
        path: 'package.json',
        file: makeFile('package.json', JSON.stringify({ name: 'my-service', dependencies: {} })),
      },
    ];

    const result = await inspectFolderFiles(entries);
    expect(result.framework).toBe('Node');
    expect(result.name).toBe('my-service');
  });
});

