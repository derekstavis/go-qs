package qs

import "testing"

func TestParse(t *testing.T) {
	Parse("foo")
	//shouldequal "foo" => nil
	Parse("foo=")
	//shouldequal "foo" => ""
	Parse("foo=bar")
	//shouldequal "foo" => "bar"
	Parse("foo=\"bar\"")
	//shouldequal "foo" => "\"bar\""

	Parse("foo=bar&foo=quux")
	//shouldequal "foo" => "quux"
	Parse("foo&foo=")
	//shouldequal "foo" => ""
	Parse("foo=1&bar=2")
	//shouldequal "foo" => "1", "bar" => "2"
	Parse("&foo=1&&bar=2")
	//shouldequal "foo" => "1", "bar" => "2"
	Parse("foo&bar=")
	//shouldequal "foo" => nil, "bar" => ""
	Parse("foo=bar&baz=")
	//shouldequal "foo" => "bar", "baz" => ""
	Parse("my+weird+field=q1%212%22%27w%245%267%2Fz8%29%3F")
	//shouldequal "my weird field" => "q1!2\"'w$5&7/z8)?"

	Parse("a=b&pid%3D1234=1023")
	//shouldequal "pid=1234" => "1023", "a" => "b"

	Parse("foo[]")
	//shouldequal "foo" => [nil]
	Parse("foo[]=")
	//shouldequal "foo" => [""]
	Parse("foo[]=bar")
	//shouldequal "foo" => ["bar"]

	Parse("foo[]=1&foo[]=2")
	//shouldequal "foo" => ["1", "2"]
	Parse("foo=bar&baz[]=1&baz[]=2&baz[]=3")
	//shouldequal "foo" => "bar", "baz" => ["1", "2", "3"]
	Parse("foo[]=bar&baz[]=1&baz[]=2&baz[]=3")
	//shouldequal "foo" => ["bar"], "baz" => ["1", "2", "3"]

	Parse("x[y][z]=1")
	//shouldequal "x" => {"y" => {"z" => "1"}}
	Parse("x[y][z][]=1")
	//shouldequal "x" => {"y" => {"z" => ["1"]}}
	Parse("x[y][z]=1&x[y][z]=2")
	//shouldequal "x" => {"y" => {"z" => "2"}}
	Parse("x[y][z][]=1&x[y][z][]=2")
	//shouldequal "x" => {"y" => {"z" => ["1", "2"]}}

	Parse("x[y][][z]=1")
	//shouldequal "x" => {"y" => [{"z" => "1"}]}
	Parse("x[y][][z][]=1")
	//shouldequal "x" => {"y" => [{"z" => ["1"]}]}
	Parse("x[y][][z]=1&x[y][][w]=2")
	//shouldequal "x" => {"y" => [{"z" => "1", "w" => "2"}]}

	Parse("x[y][][v][w]=1")
	//shouldequal "x" => {"y" => [{"v" => {"w" => "1"}}]}
	Parse("x[y][][z]=1&x[y][][v][w]=2")
	//shouldequal "x" => {"y" => [{"z" => "1", "v" => {"w" => "2"}}]}

	Parse("x[y][][z]=1&x[y][][z]=2")
	//shouldequal "x" => {"y" => [{"z" => "1"}, {"z" => "2"}]}
	Parse("x[y][][z]=1&x[y][][w]=a&x[y][][z]=2&x[y][][w]=3")
	//shouldequal "x" => {"y" => [{"z" => "1", "w" => "a"}, {"z" => "2", "w" => "3"}]}

	t.Log("---------------------------------------------------------")

	Parse("x[y]=1&x[y]z=2")
	//shouldraise(TypeError)
	//shouldequal "expected Hash (got String) for param `y'"

	Parse("x[y]=1&x[]=1")
	//shouldraise(TypeError)
	//shouldequal "expected Array (got Hash) for param `x'"

	Parse("x[y]=1&x[y][][w]=2")
	//shouldraise(TypeError)
	//shouldequal "expected Array (got String) for param `y'"
}
