import assert from 'node:assert/strict';
import { mkdtemp, mkdir, readFile, rm, writeFile, chmod } from 'node:fs/promises';
import os from 'node:os';
import path from 'node:path';
import process from 'node:process';
import test from 'node:test';
import { spawnSync } from 'node:child_process';

async function withTempPackage(fn) {
  const tempRoot = await mkdtemp(path.join(os.tmpdir(), 'smokepod-npm-run-'));
  const packageRoot = path.join(tempRoot, 'npm');
  const binDir = path.join(packageRoot, 'bin');
  const vendorDir = path.join(packageRoot, 'vendor');

  await mkdir(binDir, { recursive: true });
  await mkdir(vendorDir, { recursive: true });
  await writeFile(path.join(packageRoot, 'package.json'), JSON.stringify({ name: 'smokepod', version: '1.2.3', type: 'module' }, null, 2));
  await writeFile(path.join(packageRoot, 'install.mjs'), await readFile(new URL('../install.mjs', import.meta.url)));
  await writeFile(path.join(binDir, 'run.js'), await readFile(new URL('../bin/run.js', import.meta.url)));
  await chmod(path.join(binDir, 'run.js'), 0o755);

  try {
    await fn({ packageRoot, binDir, vendorDir });
  } finally {
    await rm(tempRoot, { recursive: true, force: true });
  }
}

function runLauncher(packageRoot, args = []) {
  return spawnSync(process.execPath, [path.join(packageRoot, 'bin', 'run.js'), ...args], {
    cwd: packageRoot,
    encoding: 'utf8'
  });
}

test('fails with a recovery message when the vendor binary is missing', async () => {
  await withTempPackage(async ({ packageRoot }) => {
    const result = runLauncher(packageRoot, ['--version']);
    assert.equal(result.status, 1);
    assert.match(result.stderr, /Re-run npm install or set SMOKEPOD_BINARY/);
  });
});

test('spawns the vendor binary, forwards args, and preserves exit codes', async () => {
  await withTempPackage(async ({ packageRoot, vendorDir }) => {
    const vendorPath = path.join(vendorDir, 'smokepod');
    await writeFile(
      vendorPath,
      [
        '#!/usr/bin/env node',
        "process.stdout.write(JSON.stringify(process.argv.slice(2)));",
        'process.exit(23);'
      ].join('\n')
    );
    await chmod(vendorPath, 0o755);

    const result = runLauncher(packageRoot, ['--json', 'config.yaml']);

    assert.equal(result.status, 23);
    assert.equal(result.stdout, '["--json","config.yaml"]');
    assert.equal(result.stderr, '');
  });
});
