	"flag"
	"log"
var nocleanup = flag.Bool("nocleanup", false, "don't clean up git state after tests are run")

	if *nocleanup {
		log.Println("directory", dir)
	} else {
		defer cleanup()
	}
	repo, err := Open(filepath.Join(dir, "repo"), "adir/", "master")
	if err != nil {
		t.Fatal(err)
	}
	commits, err := repo.Log()
	if err != nil {
		t.Fatal(err)
	}
	patch, err := repo.Patch(c.Digest)
	if err != nil {
		t.Fatal(err)
	}
func TestPatchApply(t *testing.T) {
	dir, cleanup := testutil.TempDir(t, "", "")
	if *nocleanup {
		log.Println("directory", dir)
	} else {
		defer cleanup()
	}
	shell(t, dir, `
		mkdir repos
		
		# Set up source repository and add a couple of commits:
		# - add a file to dir1
		# - move this file to dir2
		git init --bare repos/src
		git clone repos/src src
		cd src
		git config user.email you@example.com
		git config user.name "your name"
		mkdir dir1
		echo "test file" > dir1/file1
		git add dir1
		git commit -m'first commit'
		mkdir dir2
		git mv dir1/file1 dir2
		git commit -m'second commit'
		git push
		
		cd ..
		
		# Set up second, empty repository. Note that grit cannot
		# initialize empty repositories, so we add a first commit.
		git init --bare repos/dst
		git clone repos/dst dst
		cd dst
		git config user.email you@example.com
		git config user.name "your name"
		echo license > LICENSE
		git add .
		git commit -m'first commit'
		git push
	`)
	src, err := Open(filepath.Join(dir, "repos/src"), "dir2/", "master")
	if err != nil {
		t.Fatal(err)
	}
	dst, err := Open(filepath.Join(dir, "repos/dst"), "", "master")
	if err != nil {
		t.Fatal(err)
	}
	// Needs to be configured for committer.
	dst.Configure("user.email", "committer@grailbio.com")
	dst.Configure("user.name", "committer")
	commits, err := src.Log()
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(commits), 1; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	patch, err := src.Patch(commits[0].Digest)
	if err != nil {
		t.Fatal(err)
	}
	if len(patch.Diffs) == 0 {
		t.Fatal("empty patch")
	}
	if err := dst.Apply(patch); err != nil {
		t.Fatalf("failed to apply patch: %v\n%s", err, patch.Patch())
	}
	if err := dst.Push("origin", "master"); err != nil {
		t.Fatal(err)
	}
	// Make sure the file is actually there.
	shell(t, dir, `
		git -C dst pull
		cmp src/dir2/file1 dst/file1 || error file1
	`)

}

	script = `
		function error {
			echo "$@" 1>&2
			exit 1
		}
	` + script