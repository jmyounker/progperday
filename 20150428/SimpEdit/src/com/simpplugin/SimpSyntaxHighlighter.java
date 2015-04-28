package com.simpplugin;

        import com.intellij.lexer.FlexAdapter;
        import com.intellij.lexer.Lexer;
        import com.intellij.openapi.editor.SyntaxHighlighterColors;
        import com.intellij.openapi.editor.colors.TextAttributesKey;
        import com.intellij.openapi.editor.markup.TextAttributes;
        import com.intellij.openapi.fileTypes.SyntaxHighlighterBase;
        import com.intellij.psi.TokenType;
        import com.intellij.psi.tree.IElementType;
        import com.simpplugin.psi.SimpType;
        import org.jetbrains.annotations.NotNull;

        import java.awt.*;
        import java.io.Reader;

        import static com.intellij.openapi.editor.colors.TextAttributesKey.createTextAttributesKey;
/**
 * Created by jeff on 4/28/15.
 */
public class SimpSyntaxHighlighter extends SyntaxHighlighterBase {
    public static final TextAttributesKey OPERATOR = createTextAttributesKey("SIMP_SEPARATOR", SyntaxHighlighterColors.OPERATION_SIGN);
    public static final TextAttributesKey NUM = createTextAttributesKey("SIMP_NUM", SyntaxHighlighterColors.STRING);

    static final TextAttributesKey BAD_CHARACTER = createTextAttributesKey("SIMPLE_BAD_CHARACTER",
            new TextAttributes(Color.RED, null, null, null, Font.BOLD));

    private static final TextAttributesKey[] BAD_CHAR_KEYS = new TextAttributesKey[]{BAD_CHARACTER};
    private static final TextAttributesKey[] OPERATOR_KEYS = new TextAttributesKey[]{OPERATOR};
    private static final TextAttributesKey[] NUM_KEYS = new TextAttributesKey[]{NUM};
    private static final TextAttributesKey[] EMPTY_KEYS = new TextAttributesKey[0];

    @NotNull
    @Override
    public Lexer getHighlightingLexer() {
        return new FlexAdapter(new SimpLexer((Reader) null));
    }

    @NotNull
    @Override
    public TextAttributesKey[] getTokenHighlights(IElementType tokenType) {
        if (tokenType.equals(SimpType.OPERATOR)) {
            return OPERATOR_KEYS;
        } else if (tokenType.equals(SimpType.NUM)) {
            return NUM_KEYS;
        } else if (tokenType.equals(TokenType.BAD_CHARACTER)) {
            return BAD_CHAR_KEYS;
        } else {
            return EMPTY_KEYS;
        }
    }
}
