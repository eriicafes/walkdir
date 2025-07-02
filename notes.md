## Old Algo
Time complexity: O(n*m)
    where n is the `dir` and m is the files in the _n_

  The function WalkFilesWithLayout scans a given filesystem (fsys) to identify       
  and associate content files with layout files.


   1. It walks the entire directory structure, separating files into two
      categories:
       * Layout files: Special files identified by a specific layoutFilename
         (e.g., layout.html).
       * Content files: All other files that match the given file ext and are        
         within the specified dir.

   2. It then determines a hierarchy of layouts for each content file. A
      content file inherits layouts from its own directory as well as from all       
      of its parent directories.


   3. The final output is a map where each key is the path to a content file,        
      and the value is an ordered slice of strings containing the paths of all       
      applicable layouts (from the outermost directory to the innermost)
      followed by the content file's own path.


## New algo
The algorithm is as follows:
Time complexity: O(n)
To make it more efficent we use a [Trie](https://www.youtube.com/watch?v=zIjfhVPRZCg)


   1. Build a Trie (O(S)): Walk the filesystem once (where S is the total
      number of path segments, which is proportional to n). Insert every
      relevant file path into the Trie. Each node in the Trie will represent a       
      directory.
       * If a file is a layout file, we store its path in the corresponding
         directory node within the Trie.
       * If a file is a content file, we add its path to a list within its
         directory node.


   2. Traverse the Trie (O(T)): Perform a single Depth-First Search (DFS) from       
      the root of the Trie (where T is the number of nodes in the trie, also
      proportional to n).
       * As we traverse down the tree, we maintain a stack of the layouts
         encountered in the parent directories.
       * When we visit a directory node, we add its own layout (if any) to the       
         stack.
       * For every content file in that node, we now have its complete,
         ordered list of inherited layouts from the stack. We record this
         result.

  This approach processes each file and directory only once, reducing the
  complexity to be proportional to the total number of files and directories,        
  which is effectively O(n).

