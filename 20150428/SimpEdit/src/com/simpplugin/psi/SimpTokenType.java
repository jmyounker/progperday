package com.simpplugin.psi;

import com.intellij.psi.tree.IElementType;
import com.simpplugin.SimpLanguage;
import org.jetbrains.annotations.NonNls;
import org.jetbrains.annotations.NotNull;

/**
 * Created by jeff on 4/26/15.
 */
public class SimpTokenType extends IElementType {

    public SimpTokenType(@NotNull @NonNls String debugName) {
        super(debugName, SimpLanguage.INSTANCE);
        }

    @Override
    public String toString() {
        return "SimpType." + super.toString();
    }
}
