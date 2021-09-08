package consts

import "testing"

func TestTipoSecao(t *testing.T) {
	type args struct {
		tipo int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "Alcatéia", args: args{tipo: 1}, want: "Alcateia"},
		{name: "Tropa Escoteira", args: args{tipo: 2}, want: "Tropa"},
		{name: "Tropa Sênior", args: args{tipo: 3}, want: "Tropa sênior"},
		{name: "Clã Pioneiro", args: args{tipo: 4}, want: "Clã"},
		{name: "Inexistente", args: args{tipo: 0}, want: "Seção Não identificada"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TipoSecao(tt.args.tipo); got != tt.want {
				t.Errorf("TipoSecao() = %v, want %v", got, tt.want)
			}
		})
	}
}
