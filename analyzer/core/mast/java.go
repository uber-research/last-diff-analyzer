//  Copyright (c) 2023 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mast

import "fmt"

// The following constants represent legal Java modifier names.
const (
	PublicMod       = "public"
	ProtectedMod    = "protected"
	PrivateMod      = "private"
	AbstractMod     = "abstract"
	StaticMod       = "static"
	FinalMod        = "final"
	StrictfpMod     = "strictfp"
	DefaultMod      = "default"
	SynchronizedMod = "synchronized"
	NativeMod       = "native"
	TransientMod    = "transient"
	VolatileMod     = "volatile"
)

//
// Java-specific node definitions, all nodes are prefixed with Java.
//

// JavaTernaryExpression node represents a ternary expression in Java (e.g., a ? b : c).
type JavaTernaryExpression struct {
	// Condition is the condition for the expression.
	Condition Expression
	// Consequence is the result of the expression if Condition evaluates to true.
	Consequence Expression
	// Alternative is the result of the expression if Condition evaluates to false.
	Alternative Expression
}

// node implementation for JavaTernaryExpression.
func (n *JavaTernaryExpression) node() {}

// expr implementation for JavaTernaryExpression.
func (n *JavaTernaryExpression) expr() {}

// JavaAnnotatedType node represents the annotated type in Java (e.g., @NotNull int).
type JavaAnnotatedType struct {
	// Annotations is the slice of annotations for the type.
	Annotations []*Annotation
	// Type is the actual annotated type.
	Type Expression
}

// node implementation for JavaAnnotatedType.
func (n *JavaAnnotatedType) node() {}

// expr implementation for JavaAnnotatedType.
func (n *JavaAnnotatedType) expr() {}

// JavaGenericType node represents a generic type for Java (e.g., ArrayList<int> or Map<int, int>).
type JavaGenericType struct {
	// Name is the name of the type.
	Name Expression
	// Arguments is the slice of types in the arguments of the generic type.
	Arguments []Expression
}

// node implementation for JavaGenericType.
func (n *JavaGenericType) node() {}

// expr implementation for JavaGenericType.
func (n *JavaGenericType) expr() {}

// JavaWildcard node represents a wildcard type used in generic types (e.g., Map<?>).
// Note that only one of Super or Extends should be non-nil.
type JavaWildcard struct {
	// Super is the optional type for "? super type".
	Super Expression
	// Extends is the optional type for "? extends type"
	Extends Expression
}

// node implementation for JavaWildcard.
func (n *JavaWildcard) node() {}

// expr implementation for JavaWildcard.
func (n *JavaWildcard) expr() {}

// JavaArrayType node represents an array type in Java (e.g., String [][]).
type JavaArrayType struct {
	Name       Expression
	Dimensions []*JavaDimension
}

// node implementation for JavaArrayType.
func (n *JavaArrayType) node() {}

// expr implementation for JavaArrayType.
func (n *JavaArrayType) expr() {}

// JavaDimension node represents a "[]" in ArrayType.
type JavaDimension struct {
	// Length is an optional length of the dimension.
	Length Expression
	// Annotations is the optional slice of annotations for the current dimension.
	Annotations []*Annotation
}

// node implementation for JavaDimension.
func (n *JavaDimension) node() {}

// expr implementation for JavaDimension.
func (n *JavaDimension) expr() {}

// JavaInstanceOfExpression node represents an instanceof expression in Java (e.g., "foo instanceof Bar").
type JavaInstanceOfExpression struct {
	// Operand is the operand of the instanceof operation.
	Operand Expression
	// Type is the type to check for the instanceof operation.
	Type Expression
}

// node implementation for JavaInstanceOfExpression.
func (n *JavaInstanceOfExpression) node() {}

// expr implementation for JavaInstanceOfExpression.
func (n *JavaInstanceOfExpression) expr() {}

// JavaLiteralModifier node represents a literal modifier expression in java (e.g. "public", "protected", etc.).
type JavaLiteralModifier struct {
	// Modifier is the literal modifier.
	Modifier string
}

// node implementation for JavaLiteralModifier.
func (n *JavaLiteralModifier) node() {}

// expr implementation for JavaLiteralModifier.
func (n *JavaLiteralModifier) expr() {}

// JavaTryStatement node represents a try-catch-finally statement in Java
// (e.g., "try {...} catch (Exception e) {...} finally {...}").
type JavaTryStatement struct {
	// Resources is an optional list of resources for a try-with-resources statement(e.g., "try (file) {...}").
	// Each resource has one of the following forms:
	// (1) a plain Identifier or AccessPath;
	// (2) a "mini" variable declaration (Modifiers, Type, Name, Dimensions, Value).
	// Therefore, for case (1) the element will be an Identifier or AccessPath node wrapped in a ExpressionStatement
	// node and for case (2) it will be a VariableDeclaration node wrapped in a DeclarationStatement.
	Resources []Statement
	// Body is the Block node for the body of the try statement.
	Body *Block
	// CatchClauses is the slice of catch clauses for the try statement.
	CatchClauses []*JavaCatchClause
	// Finally is the finally clause for the try statement, nil if there is none.
	Finally *JavaFinallyClause
}

// node implementation for JavaTryStatement.
func (n *JavaTryStatement) node() {}

// stmt implementation for JavaTryStatement.
func (n *JavaTryStatement) stmt() {}

// JavaCatchClause node represents a catch clause for try-catch statement in Java (e.g., "catch (Exception a) {...}").
type JavaCatchClause struct {
	// Parameter is a JavaCatchFormalParameter node that wraps the formal parameter of the catch clause.
	Parameter *JavaCatchFormalParameter
	// Body is the block body of the catch clause
	Body *Block
}

// node implementation for JavaCatchClause.
func (n *JavaCatchClause) node() {}

// expr implementation for JavaCatchClause.
func (n *JavaCatchClause) expr() {}

// JavaCatchFormalParameter node represents the formal parameters in the JavaCatchClause node.
type JavaCatchFormalParameter struct {
	// Modifiers is the slice of modifiers for the exception types, each element could either be a JavaLiteralModifier
	// or an Annotation.
	Modifiers []Expression
	// Types is the slice of exception types for the catch clause, connected by "|".
	Types []Expression
	// Name is the name of the exception.
	Name *Identifier
	// Dimensions is the optional slice of dimensions for the name.
	Dimensions []*JavaDimension
}

// node implementation for JavaCatchFormalParameter.
func (n *JavaCatchFormalParameter) node() {}

// expr implementation for JavaCatchFormalParameter.
func (n *JavaCatchFormalParameter) expr() {}

// JavaFinallyClause node represents a finally clause in try statement in Java (e.g., "finally {...}").
type JavaFinallyClause struct {
	// Body is the Block node for the body of the finally clause.
	Body *Block
}

// node implementation for JavaFinallyClause.
func (n *JavaFinallyClause) node() {}

// expr implementation for JavaFinallyClause.
func (n *JavaFinallyClause) expr() {}

// JavaWhileStatement node represents a while statement in Java (e.g., "while a { foo(); }")
type JavaWhileStatement struct {
	// Condition is the condition for the while statement.
	Condition Expression
	// Body is the body of the while statement, it could be a Block node for a block of statements.
	Body Statement
}

// node implementation for JavaWhileStatement.
func (n *JavaWhileStatement) node() {}

// expr implementation for JavaWhileStatement.
func (n *JavaWhileStatement) stmt() {}

// JavaThrowStatement node represents a throw statement in Java (e.g., "throw ex;")
type JavaThrowStatement struct {
	// Expr is the expression for the throw statement.
	Expr Expression
}

// node implementation for JavaThrowStatement.
func (n *JavaThrowStatement) node() {}

// expr implementation for JavaThrowStatement.
func (n *JavaThrowStatement) stmt() {}

// JavaAssertStatement node represents an assert statement in Java (e.g., "assert a == b;" or "assert a == b : str;")
type JavaAssertStatement struct {
	// Condition is the expression for the assert statement.
	Condition Expression
	// ErrorString is the error string that is used to generate the AssertionError.
	ErrorString Expression
}

// node implementation for JavaAssertStatement.
func (n *JavaAssertStatement) node() {}

// expr implementation for JavaAssertStatement.
func (n *JavaAssertStatement) stmt() {}

// JavaSynchronizedStatement node represents a synchronized statement in Java (e.g., "synchronized (a) {a++;}")
type JavaSynchronizedStatement struct {
	// Expr is the expression for the synchronized statement.
	Expr Expression
	// Body is the block for the synchronized statement.
	Body *Block
}

// node implementation for JavaSynchronizedStatement.
func (n *JavaSynchronizedStatement) node() {}

// expr implementation for JavaSynchronizedStatement.
func (n *JavaSynchronizedStatement) stmt() {}

// JavaDoStatement node represents a do-while statement in Java (e.g., "do {...} while (c);")
type JavaDoStatement struct {
	// Body is the body of the while statement, it could be a Block node for a block of statements.
	Body Statement
	// Condition is the condition for the do statement.
	Condition Expression
}

// node implementation for JavaDoStatement.
func (n *JavaDoStatement) node() {}

// expr implementation for JavaDoStatement.
func (n *JavaDoStatement) stmt() {}

// JavaVariableDeclarator node represents part of variable declaration with an optional initializer but without type
// information (e.g. "a[][]", "a = 5"). Note that tree-sitter sometimes (but not always) generates variable_declarator
// (which will be translated to this node) inside a few nodes (e.g., spread_parameter). This node is only meant to be
// an intermediate node that will be unwrapped by other nodes, in other words, this node will _never_ appear in the
// final MAST. Therefore we deliberately omit implementing a method that would put this node into one of the AST node
// categories (e.g. stmt() method for statements, or expr() method for expressions). Since all nodes in MAST fall in
// one of the categories, this disallows the appearance of this node in the final MAST.
type JavaVariableDeclarator struct {
	// Name is the name of the variable.
	Name *Identifier
	// Dimensions is an optional slice of dimensions for the variable.
	Dimensions []*JavaDimension
	// Initializer is an optional initializer for the variable.
	Initializer Expression
}

// node implementation for JavaVariableDeclarator.
func (n *JavaVariableDeclarator) node() {}

// JavaEnhancedForStatement node represents an enhanced for statement in Java (e.g., "for (String x: lst) {...}").
type JavaEnhancedForStatement struct {
	// Modifiers is an optional slice of modifiers for the for statement, the element could either be a
	// *mast.Annotation node or *mast.JavaLiteralModifier.
	Modifiers []Expression
	// Type is the type of the element in the iterable.
	Type Expression
	// Name is the name of the element in the iterable.
	Name *Identifier
	// Dimensions is an optional slice of dimensions for the element in the iterable.
	Dimensions []*JavaDimension
	// Iterable is the expression of the iterable.
	Iterable Expression
	// Body is the body of the for statement, this could either be a a single statement or a *mast.Block for a block of
	// statements.
	Body Statement
}

// node implementation for JavaEnhancedForStatement.
func (n *JavaEnhancedForStatement) node() {}

// expr implementation for JavaEnhancedForStatement.
func (n *JavaEnhancedForStatement) stmt() {}

// JavaModuleDeclaration node represents a module declaration in Java (e.g., "module A {...}").
type JavaModuleDeclaration struct {
	// Annotations is the optional slice of annotations for the module declaration.
	Annotations []*Annotation
	// IsOpen indicates if module is declared with "open" keyword.
	IsOpen bool
	// Name is the name of the module, could either be an Identifier node or an AccessPath node.
	Name Expression
	// Directives is the slice of module directives.
	Directives []*JavaModuleDirective
}

// node implementation for JavaModuleDeclaration.
func (n *JavaModuleDeclaration) node() {}

// decl implementation for JavaModuleDeclaration.
func (n *JavaModuleDeclaration) decl() {}

// JavaModuleDirective node represents a directive in the module declaration (e.g., "requires a;").
type JavaModuleDirective struct {
	// Keyword is the keyword for the directive, e.g., "requires" and "exports".
	Keyword string
	// Exprs is the slice of expressions for the directive. Depending on the keyword, each element in the expressions
	// has different meaning. See JLS [1] for further details on this.
	// [1] https://docs.oracle.com/javase/specs/jls/se11/html/jls-7.html#jls-7.7
	Exprs []Expression
}

// node implementation for JavaModuleDirective.
func (n *JavaModuleDirective) node() {}

// decl implementation for JavaModuleDirective.
func (n *JavaModuleDirective) decl() {}

// JavaTypeParameter node represents a generic type parameter in class, interface or method declaration (e.g., "class Test<A, B>").
// This node is only intended to be a direct field of other nodes, therefore we deliberately omit implementing a method
// that would put this node into one of the AST node categories (e.g. stmt() method for statements, or expr() method for
// expressions).
type JavaTypeParameter struct {
	// Annotations is an optional slice of annotations for the type parameter.
	Annotations []*Annotation
	// Type is the type for the type parameter.
	Type Expression
	// Extends is an optional slice of types that the current type extends.
	Extends []Expression
}

// node implementation for JavaTypeParameter.
func (n *JavaTypeParameter) node() {}

// JavaClassDeclaration node represents a class declaration in Java (e.g., "class Test {...}").
type JavaClassDeclaration struct {
	// Modifiers is an optional slice of modifiers for the class.
	Modifiers []Expression
	// Name is the name of the class.
	Name *Identifier
	// TypeParameters is an optional slice of JavaTypeParameters for the class (e.g. "class Test<A, B>").
	TypeParameters []*JavaTypeParameter
	// SuperClass is the optional super class of the class.
	SuperClass Expression
	// Interfaces is the optional slice of interfaces this class implements.
	Interfaces []Expression
	// Body is the body of the class declaration.
	Body []Declaration
}

// node implementation for JavaClassDeclaration.
func (n *JavaClassDeclaration) node() {}

// decl implementation for JavaClassDeclaration.
func (n *JavaClassDeclaration) decl() {}

// JavaInterfaceDeclaration node represents an interface declaration in Java (e.g., "interface I {..}").
type JavaInterfaceDeclaration struct {
	// Modifiers is an optional slice of modifiers for the interface.
	Modifiers []Expression
	// Name is the name of the interface.
	Name *Identifier
	// TypeParameters is an optional slice of JavaTypeParameters for the interface (e.g. "interface Test<A, B>").
	TypeParameters []*JavaTypeParameter
	// Extends is an optional slice of interfaces this interface extends.
	Extends []Expression
	// Body is the body of the interface.
	Body []Declaration
}

// node implementation for JavaInterfaceDeclaration.
func (n *JavaInterfaceDeclaration) node() {}

// decl implementation for JavaInterfaceDeclaration.
func (n *JavaInterfaceDeclaration) decl() {}

// JavaEnumDeclaration node represents an enum declaration in Java (e.g., "enum Color {...}").
type JavaEnumDeclaration struct {
	// Modifiers is an optional slice of modifiers for the enum declaration.
	Modifiers []Expression
	// Name is the name of the enum.
	Name *Identifier
	// Interfaces is the optional slice of interfaces this enum implements.
	Interfaces []Expression
	// Body is the body of the enum.
	Body []Declaration
}

// node implementation for JavaEnumDeclaration.
func (n *JavaEnumDeclaration) node() {}

// decl implementation for JavaEnumDeclaration.
func (n *JavaEnumDeclaration) decl() {}

// JavaEnumConstantDeclaration node represents an enum constant inside the JavaEnumDeclaration (e.g., "enum T {A, B}").
type JavaEnumConstantDeclaration struct {
	// Modifiers is an optional slice of modifiers for the enum constant.
	Modifiers []Expression
	// Name is the name of the enum constant.
	Name *Identifier
	// Arguments is an optional list of arguments for the enum constant.
	Arguments []Expression
	// Body is an optional body of the enum constant.
	Body []Declaration
}

// node implementation for JavaEnumConstantDeclaration.
func (n *JavaEnumConstantDeclaration) node() {}

// decl implementation for JavaEnumConstantDeclaration.
func (n *JavaEnumConstantDeclaration) decl() {}

// JavaClassInitializer node represents an initializer in class declaration body in Java (e.g.,
// "class T { static {...} }" or "class T { {...} }").
type JavaClassInitializer struct {
	// IsStatic indicates whether the initializer is a static initializer.
	IsStatic bool
	// Block is the block for the static initializer.
	Block *Block
}

// node implementation for JavaClassInitializer.
func (n *JavaClassInitializer) node() {}

// decl implementation for JavaClassInitializer.
func (n *JavaClassInitializer) decl() {}

// JavaAnnotationDeclaration node represents an annotation declaration in Java (e.g., "@interface Test {...}").
type JavaAnnotationDeclaration struct {
	// Modifiers is an optional list of modifiers for the annotation declaration.
	Modifiers []Expression
	// Name is the name of the annotation.
	Name *Identifier
	// Body is the body of the annotation declaration.
	Body []Declaration
}

// node implementation for JavaAnnotationDeclaration.
func (n *JavaAnnotationDeclaration) node() {}

// decl implementation for JavaAnnotationDeclaration.
func (n *JavaAnnotationDeclaration) decl() {}

// JavaAnnotationElementDeclaration node represents an element declaration inside annotation declaration in Java (
// e.g., "int A() default 2;").
type JavaAnnotationElementDeclaration struct {
	// Modifiers is an optional list of modifiers for the annotation element declaration.
	Modifiers []Expression
	// Type is the type of the element.
	Type Expression
	// Name is the name of the element.
	Name *Identifier
	// Value is an optional default value for the element.
	Value Expression
	// Dimensions is an optional list of dimensions for the element.
	Dimensions []*JavaDimension
}

// node implementation for JavaAnnotationElementDeclaration.
func (n *JavaAnnotationElementDeclaration) node() {}

// decl implementation for JavaAnnotationElementDeclaration.
func (n *JavaAnnotationElementDeclaration) decl() {}

// JavaMethodReference node represents a method reference expression in Java (e.g., "SomeClass::someMethod")
type JavaMethodReference struct {
	// Class is the class name for the method reference expression.
	Class Expression
	// TypeArguments is an optional list of type parameters for the method.
	TypeArguments []Expression
	// Method is the method name of the method reference expression.
	Method *Identifier
}

// node implementation for JavaMethodReference.
func (n *JavaMethodReference) node() {}

// expr implementation for JavaMethodReference.
func (n *JavaMethodReference) expr() {}

// JavaClassLiteral node represents a class literal expression in Java (e.g., "String[].class"). Note that class
// literals without dimensions such as "String.class" will be classified as AccessPath nodes.
type JavaClassLiteral struct {
	// Type is the type of the class literal.
	Type Expression
}

// node implementation for JavaClassLiteral.
func (n *JavaClassLiteral) node() {}

// expr implementation for JavaClassLiteral.
func (n *JavaClassLiteral) expr() {}

//
// Definitions for language-specific fields that extend the generic MAST nodes.
//

// JavaFunctionDeclarationFields node stores the Java-specific fields for the generic FunctionDeclaration node, to
// represent a method declaration in Java (e.g., "public void test() {...}").
type JavaFunctionDeclarationFields struct {
	// Modifiers is an optional list of modifiers (either an *Annotation or *JavaLiteralModifier).
	Modifiers []Expression
	// TypeParameters is an optional list of type parameters for the method declaration.
	TypeParameters []*JavaTypeParameter
	// Annotations is an optional list of annotations for the method declaration.
	Annotations []*Annotation
	// Dimension is an optional list of dimensions for the method declaration.
	Dimensions []*JavaDimension
	// Throws is an optional list of exception types that this method can throw.
	Throws []Expression
}

// node implementation for JavaFunctionDeclarationFields.
func (n *JavaFunctionDeclarationFields) node() {}

// langFunctionDeclarationFields implementation for JavaFunctionDeclarationFields.
func (n *JavaFunctionDeclarationFields) langFunctionDeclarationFields() {}

// JavaParameterDeclarationFields node stores the Java-specific fields for the generic ParameterDeclaration node, to
// represent a parameter declaration in Java (e.g., "void test(@NotNull int a[]) {...}").
type JavaParameterDeclarationFields struct {
	// IsReceiver indicates wither this parameter declaration is a receiver parameter, a special parameter in Java
	// (e.g., "(@Test T name. this)").
	IsReceiver bool
	// Modifiers is an optional slice of modifiers, only available for Java.
	Modifiers []Expression
	// Dimensions is an optional slice of dimensions, only available for Java.
	Dimensions []*JavaDimension
}

// node implementation for JavaParameterDeclarationFields.
func (n *JavaParameterDeclarationFields) node() {}

// langParameterDeclarationFields implementation for JavaParameterDeclarationFields.
func (n *JavaParameterDeclarationFields) langParameterDeclarationFields() {}

// JavaEntityCreationExpressionFields node stores the Java-specific
// fields for the generic EntityCreationExpression node, to represent
// optional dimensions and optional body of class declaration.
type JavaEntityCreationExpressionFields struct {
	// Dimensions is an optional slice of dimensions.
	Dimensions []*JavaDimension
	// Body is the body of the class declaration.
	Body []Declaration
}

// node implementation for JavaEntityCreationExpressionFields.
func (n *JavaEntityCreationExpressionFields) node() {}

// langEntityCreationExpressionFields implementation for JavaEntityCreationExpressionFields.
func (n *JavaEntityCreationExpressionFields) langEntityCreationExpressionFields() {}

// JavaCallExpressionFields node stores the Java-specific fields for the generic
// CallExpression node, to represent optional type arguments for:
// (1) Method Invocations, e.g.,
//
//	"pkg.Clazz.<String, Integer>foo("HelloWorld!", 42);"
//
// (2) Explicit construct call, e.g.,
//
//	(a) this(arguments);
//	(b) super(arguments);
//	(c) <A, B...>this(arguments);
//	(d) <A, B...>super(arguments);
//	(e) Expression.<A, B...>super(arguments);
//	See https://docs.oracle.com/javase/specs/jls/se11/html/jls-8.html#jls-8.8.7.1
//	for detailed explanations.
type JavaCallExpressionFields struct {
	// TypeArguments is an optional list of type arguments for the method call
	// or constructor call.
	TypeArguments []Expression
}

// node implementation for JavaCallExpressionFields.
func (n *JavaCallExpressionFields) node() {}

// langCallExpressionFields implementation for JavaCallExpressionFields.
func (n *JavaCallExpressionFields) langCallExpressionFields() {}

// JavaVariableDeclarationFields node stores the Java-specific
// fields for the generic VariableDeclaration node, to represent
// optional annotations and optional dimensions.
type JavaVariableDeclarationFields struct {
	// Modifiers is an optional list of modifiers (either an
	// *Annotation or *JavaLiteralModifier).
	Modifiers []Expression
	// Dimensions is an optional slice of dimensions.
	Dimensions []*JavaDimension
}

// node implementation for JavaVariableDeclarationFields.
func (n *JavaVariableDeclarationFields) node() {}

// langVariableDeclarationFields implementation for JavaVariableDeclarationFields.
func (n *JavaVariableDeclarationFields) langVariableDeclarationFields() {}

// SetJavaExprTypeKinds sets kinds of identifiers in a Java expression to the
// type kind
func SetJavaExprTypeKinds(expr Expression) error {
	switch e := expr.(type) {
	case *Identifier:
		// unqualified type name (e.g., Clazz)
		e.Kind = Typ
	case *AccessPath:
		// It's either a qualified type name (e.g., pkg.com.Clazz) or
		// an access path representing inner class chaing (e.g.,
		// OuterClazz.InnerClazz). In general we have no way of
		// knowing which one it is, but since at this point we only
		// handle the first identifier on the access path, this is the
		// only identifier we have to safely handle.
		//
		// The first identifier could be a (sub)package name or a
		// class name. According to the JLS
		// (https://docs.oracle.com/javase/specs/jls/se7/html/jls-6.html#jls-6.5.4.1)
		// if we have to decide between package name and type name, we
		// first look if the type declaration with the same name
		// exists (in which case, the identifier is of type kind and
		// is resolved to refer to this declaration). If such type
		// declaration does not exist, we look for package matching a
		// given identifier (and resolved it to refer to the package).
		//
		// In our case then, it is safe to assing type kind to the
		// identifier right away (and assignment of kinds to other
		// identifiers on the access path does not matter at this
		// point). If the matching type declaration is found, we
		// resolve this identifier to this type, otherwise it will
		// always remain unresolved as we do not even attempt
		// symbolication of package names.
		//
		// Let's do custom iteration instead of calling
		// extractAccessPath to make sure that we only have identifers
		// on this path (other expressions would make no sense in a
		// type name).
		e.Field.Kind = Typ
		var current Expression = e.Operand
	outer:
		for {
			switch n := current.(type) {
			case *AccessPath:
				// "recurse" into the inner parts of access path
				n.Field.Kind = Typ
				current = n.Operand
			case *Identifier:
				// end of "recursion"
				n.Kind = Typ
				break outer
			default:
				return fmt.Errorf("node %T has unexpected type when setting type kind for access path", n)
			}
		}
	case *JavaArrayType, *JavaGenericType, *JavaAnnotatedType:
		// these cases are handled during translation into respective Java nodes
	default:
		// we want to make sure that we are processing only
		// expressions of a certain shape in this function
		return fmt.Errorf("node %T has unexpected type when setting type kind", e)
	}
	return nil
}

// SetJavaCallKind sets the kind of the standalone identifier
// representing Java method call name, such as the kind of foo in
// foo().
//
// This method does some additional verification to make sure that we
// do not miss setting the kind for other relevant cases, such as for
// calls via an access path, for example a.b.foo().
func SetJavaCallKind(expr Expression, kind NameKind) error {
	switch e := expr.(type) {
	case *JavaGenericType:
		// recurse on the expression representing type name: new java.util.HashSet<>()
		return SetJavaCallKind(e.Name, kind)
	}
	// generic case
	return SetCallKind(expr, kind)
}
