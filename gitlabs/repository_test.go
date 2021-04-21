package gitlabs

import "testing"

func Test_getBranchInfo(t *testing.T) {
	tests := []struct {
		name  string
		s     string
		b     string
		want  bool
		want1 bool
		want2 string
	}{
		// TODO: Add test cases.
		{name: "1", s: `  dev     0f2e21e [origin/dev：领先 1] Merge branch 'dev'
		dev1994 0f2e21e [origin/dev：领先 1] Merge branch 'dev'
	* dev2000 0f2e21e Merge branch 'dev'
		master  0f2e21e [origin/master] Merge branch 'dev'`, b: "dev", want: true, want1: false, want2: "dev"},
		{name: "1", s: `  dev     0f2e21e [origin/dev：领先 1] Merge branch 'dev'
		dev1994 0f2e21e [origin/dev：领先 1] Merge branch 'dev'
	* dev2000 0f2e21e Merge branch 'dev'
		master  0f2e21e [origin/master] Merge branch 'dev'`, b: "dev1994", want: true, want1: false, want2: "dev"},
		{name: "1", s: `  dev     0f2e21e [origin/dev：领先 1] Merge branch 'dev'
		dev1994 0f2e21e [origin/dev：领先 1] Merge branch 'dev'
	* dev2000 0f2e21e Merge branch 'dev'
		master  0f2e21e [origin/master] Merge branch 'dev'`, b: "dev2000", want: true, want1: true, want2: ""},
		{name: "1", s: `  dev     0f2e21e [origin/dev：领先 1] Merge branch 'dev'
		dev1994 0f2e21e [origin/dev：领先 1] Merge branch 'dev'
	* dev2000 0f2e21e Merge branch 'dev'
		master  0f2e21e [origin/master] Merge branch 'dev'`, b: "master", want: true, want1: false, want2: "master"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := getBranchInfo(tt.s, tt.b)
			if got != tt.want {
				t.Errorf("getBranchInfo() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("getBranchInfo() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("getBranchInfo() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}
