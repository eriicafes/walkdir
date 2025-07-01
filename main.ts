import { walkFilesWithLayout } from "./files";

console.log("start ts");
const filesWithLayouts = walkFilesWithLayout(
    "example",
    "html",
    "layout",
    "app",
);

for (const [file, paths] of Object.entries(filesWithLayouts)) {
  console.log(
    file.slice("basic".length),
    "->",
    paths.map((p) => p.slice("basic".length))
  );
}
console.log("end ts");
