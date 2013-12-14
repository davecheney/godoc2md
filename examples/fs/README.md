
	
		
		
Package fs provides filesystem-related functions.


		

		
		<div id="pkg-examples">
			<h4>Examples</h4>
			<dl>
			
			<dd><a class="exampleLink" href="#example_Walker">Walker</a></dd>
			
			</dl>
		</div>
		







## type FileSystem
<pre>type FileSystem interface {

    <span class="comment">// ReadDir reads the directory named by dirname and returns a</span>
    <span class="comment">// list of directory entries.</span>
    ReadDir(dirname string) ([]os.FileInfo, error)

    <span class="comment">// Lstat returns a FileInfo describing the named file. If the file is a</span>
    <span class="comment">// symbolic link, the returned FileInfo describes the symbolic link. Lstat</span>
    <span class="comment">// makes no attempt to follow the link.</span>
    Lstat(name string) (os.FileInfo, error)

    <span class="comment">// Join joins any number of path elements into a single path, adding a</span>
    <span class="comment">// separator if necessary. The result is Cleaned; in particular, all</span>
    <span class="comment">// empty strings are ignored.</span>
    <span class="comment">//</span>
    <span class="comment">// The separator is FileSystem specific.</span>
    Join(elem ...string) string
}</pre>

FileSystem defines the methods of an abstract filesystem.















## type Walker
<pre>type Walker struct {
    <span class="comment">// contains filtered or unexported fields</span>
}</pre>

Walker provides a convenient interface for iterating over the
descendants of a filesystem path.
Successive calls to the Step method will step through each
file or directory in the tree, including the root. The files
are walked in lexical order, which makes the output deterministic
but means that for very large directories Walker can be inefficient.
Walker does not follow symbolic links.











### func Walk
<pre>func Walk(root string) *Walker</pre>

Walk returns a new Walker rooted at root.





### func WalkFS
<pre>func WalkFS(root string, fs FileSystem) *Walker</pre>

WalkFS returns a new Walker rooted at root on the FileSystem fs.







### func (*Walker) Err
<pre>func (w *Walker) Err() error</pre>
<p>
Err returns the error, if any, for the most recent attempt
by Step to visit a file or directory. If a directory has
an error, w will not descend into that directory.
</p>





### func (*Walker) Path
<pre>func (w *Walker) Path() string</pre>
<p>
Path returns the path to the most recent file or directory
visited by a call to Step. It contains the argument to Walk
as a prefix; that is, if Walk is called with &#34;dir&#34;, which is
a directory containing the file &#34;a&#34;, Path will return &#34;dir/a&#34;.
</p>





### func (*Walker) SkipDir
<pre>func (w *Walker) SkipDir()</pre>
<p>
SkipDir causes the currently visited directory to be skipped.
If w is not on a directory, SkipDir has no effect.
</p>





### func (*Walker) Stat
<pre>func (w *Walker) Stat() os.FileInfo</pre>
<p>
Stat returns info for the most recent file or directory
visited by a call to Step.
</p>





### func (*Walker) Step
<pre>func (w *Walker) Step() bool</pre>
<p>
Step advances the Walker to the next file or directory,
which will then be available through the Path, Stat,
and Err methods.
It returns false when the walk stops at the end of the tree.
</p>










