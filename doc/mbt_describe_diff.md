## mbt describe diff

Describe the modules in the diff between from and to commits

### Synopsis


Describe the modules in the diff between from and to commits

Works out the merge base for from and to commits and
displays all modules which have been changed between merge base and
the to commit.


```
mbt describe diff --from <commit> --to <commit> [flags]
```

### Options

```
      --from string   from commit
  -h, --help          help for diff
      --to string     to commit
```

### Options inherited from parent commands

```
      --debug       enable debugging
      --graph       format output as dot graph
      --in string   path to repo
      --json        format output as json
```

### SEE ALSO
* [mbt describe](mbt_describe.md)	 - Describe the modules in the repo

###### Auto generated by spf13/cobra on 3-Apr-2018
