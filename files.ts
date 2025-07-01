import * as fs from "node:fs";
import * as path from "path";

/**
 * walkFilesWithLayout walks a directory and for each file that matches the given extension,
 * returns a map where the key is the filename (without extension) and the value is an array
 * of layout filenames (without extension) in parent directories followed by the matched file itself.
 *
 * @param rootDir Root directory to begin walking
 * @param ext File extension to filter for (e.g. 'md')
 * @param layoutFilename Name of layout file (e.g. 'layout')
 * @param dir Optional path prefix filter (e.g. '.', or 'content')
 * @returns Map of file (without extension) to list of layouts + itself
 */
export function walkFilesWithLayout(
  rootDir: string,
  ext: string,
  layoutFilename: string,
  dir: string = "."
): Record<string, string[]> {
  const reldir = path.join(rootDir, dir);

  const groups: Record<string, string[]> = {};
  const layouts: string[] = [];

  walk(rootDir, (filePath) => {
    // skip if path does not have ext
    if (!filePath.endsWith("." + ext)) return;
    const pathWithoutExt = filePath.split("." + ext)[0];
    const filename = path.basename(pathWithoutExt);

    if (filename === layoutFilename) {
      layouts.push(pathWithoutExt);
    } else if (dir === "." || pathWithoutExt.startsWith(reldir)) {
      groups[pathWithoutExt] = [pathWithoutExt];
    }
  });

  if (layouts.length === 0) return groups;
  // sort by shortest layout first (shorter layout is parent of longer layout)
  layouts.sort((a, b) => a.length - b.length);

  for (const name of Object.keys(groups)) {
    const files: string[] = [];
    const fileDir = path.dirname(name);

    for (const layout of layouts) {
      const layoutDir = path.dirname(layout);
      if (fileDir.startsWith(layoutDir)) {
        files.push(layout);
      }
      // no need to check deeper layout files
      if (layoutDir === fileDir) break;
    }

    groups[name] = [...files, ...groups[name]];
  }

  return groups;
}

function walk(
  currentDir: string,
  walkFn: (path: string, entry: fs.Dirent) => void
) {
  const entries = fs.readdirSync(currentDir, { withFileTypes: true });

  for (const entry of entries) {
    const fullPath = path.join(currentDir, entry.name);
    if (entry.isDirectory()) walk(fullPath, walkFn);
    else walkFn(fullPath, entry);
  }
}
