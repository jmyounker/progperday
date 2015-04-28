package com.simpplugin;

import com.intellij.lexer.FlexLexer;
import com.intellij.psi.tree.IElementType;
import com.simpplugin.psi.SimpType;
import com.intellij.psi.TokenType;

%%

%class SimpLexer
%implements FlexLexer
%unicode
%function advance
%type IElementType
%eof{  return;
%eof}

CRLF= \n|\r|\r\n
WHITE_SPACE=[\ \t\f]
NUM="-"?[0-9][0-9]*
OPERATOR="+"|"-"|"*"|"/"
OPEN_PAREN="("
CLOSE_PAREN=")"
// COMMENT='//'

%state WAITING_VALUE

%%


// {COMMENT}                                    { yybegin(YYINITIAL); return SimpType.COMMENT; }

{NUM}                                        { yybegin(YYINITIAL); return SimpType.NUM; }

{OPEN_PAREN}                                 { yybegin(YYINITIAL); return SimpType.OPEN_PAREN; }

{OPERATOR}                                   { yybegin(YYINITIAL); return SimpType.OPERATOR; }

{CLOSE_PAREN}                                { yybegin(YYINITIAL); return SimpType.CLOSE_PAREN; }

{CRLF}                                       { yybegin(YYINITIAL); return SimpType.CRLF; }

{WHITE_SPACE}+                               { yybegin(YYINITIAL); return TokenType.WHITE_SPACE; }

.                                            { return TokenType.BAD_CHARACTER; }

