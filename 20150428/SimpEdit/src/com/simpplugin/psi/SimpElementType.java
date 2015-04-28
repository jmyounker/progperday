package com.simpplugin.psi;

/**
 * Created by jeff on 4/28/15.
 */

import com.intellij.psi.tree.IElementType;
import com.simpplugin.SimpLanguage;
import org.jetbrains.annotations.NonNls;
import org.jetbrains.annotations.NotNull;

public class SimpElementType extends IElementType {
    public SimpElementType(@NotNull @NonNls String debugName) {
        super(debugName, SimpLanguage.INSTANCE);
    }
}
