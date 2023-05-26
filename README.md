# aldana
Lexer/Parser/Transpiler tools

````go
type IProvider[T any] interface {
    Get() bool
    Set(b bool)
    Do(f T) error
}

type Foo struct {}

type Provider implemention IProvider[Foo] {
    flag bool

    func Get() bool {
        return this.flag
    }

    func Set(b bool) bool {
        this.flag = b
    }

    func Do(f Foo) error {
        return nil
    }
}

````