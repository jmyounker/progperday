/**
 * Created by jeff on 4/27/15.
 */
package com.simpplugin;

import com.intellij.lang.ASTNode;
import com.intellij.lang.Language;
import com.intellij.lang.ParserDefinition;
import com.intellij.lang.PsiParser;
import com.intellij.lexer.FlexAdapter;
import com.intellij.lexer.Lexer;
import com.intellij.openapi.project.Project;
import com.intellij.psi.FileViewProvider;
import com.intellij.psi.PsiElement;
import com.intellij.psi.PsiFile;
import com.intellij.psi.TokenType;
import com.intellij.psi.tree.IFileElementType;
import com.intellij.psi.tree.TokenSet;
import com.simpplugin.parser.SimpParser;
import com.simpplugin.psi.SimpFile;
import com.simpplugin.psi.SimpType;
import org.jetbrains.annotations.NotNull;

import java.io.Reader;

public class SimpParserDefinition implements ParserDefinition{
    public static final TokenSet WHITE_SPACES = TokenSet.create(TokenType.WHITE_SPACE);
    public static final TokenSet COMMENTS = TokenSet.create(TokenType.WHITE_SPACE); // NO COMMENT!
    public static final TokenSet NUM = TokenSet.create(SimpType.NUM);

    public static final IFileElementType FILE = new IFileElementType(Language.<SimpLanguage>findInstance(SimpLanguage.class));

    @NotNull
    @Override
    public Lexer createLexer(Project project) {
        return new FlexAdapter(new SimpLexer((Reader) null));
    }

    @NotNull
    public TokenSet getWhitespaceTokens() {
        return WHITE_SPACES;
    }

    @NotNull
    public TokenSet getCommentTokens() {
        return COMMENTS;
    }

    @NotNull
    @Override
    public TokenSet getStringLiteralElements() {
        return NUM;
    }

    @NotNull
    public PsiParser createParser(final Project project) {
        return new SimpParser();
    }

    @Override
    public IFileElementType getFileNodeType() {
        return FILE;
    }

    public PsiFile createFile(FileViewProvider viewProvider) {
        return new SimpFile(viewProvider);
    }

    public SpaceRequirements spaceExistanceTypeBetweenTokens(ASTNode left, ASTNode right) {
        return SpaceRequirements.MAY;
    }

    @NotNull
    public PsiElement createElement(ASTNode node) {
        return SimpType.Factory.createElement(node);
    }
}
