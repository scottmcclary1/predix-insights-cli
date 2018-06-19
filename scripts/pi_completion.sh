# bash completion for pi                                   -*- shell-script -*-

__pi_debug()
{
    if [[ -n ${BASH_COMP_DEBUG_FILE} ]]; then
        echo "$*" >> "${BASH_COMP_DEBUG_FILE}"
    fi
}

# Homebrew on Macs have version 1.3 of bash-completion which doesn't include
# _init_completion. This is a very minimal version of that function.
__pi_init_completion()
{
    COMPREPLY=()
    _get_comp_words_by_ref "$@" cur prev words cword
}

__pi_index_of_word()
{
    local w word=$1
    shift
    index=0
    for w in "$@"; do
        [[ $w = "$word" ]] && return
        index=$((index+1))
    done
    index=-1
}

__pi_contains_word()
{
    local w word=$1; shift
    for w in "$@"; do
        [[ $w = "$word" ]] && return
    done
    return 1
}

__pi_handle_reply()
{
    __pi_debug "${FUNCNAME[0]}"
    case $cur in
        -*)
            if [[ $(type -t compopt) = "builtin" ]]; then
                compopt -o nospace
            fi
            local allflags
            if [ ${#must_have_one_flag[@]} -ne 0 ]; then
                allflags=("${must_have_one_flag[@]}")
            else
                allflags=("${flags[*]} ${two_word_flags[*]}")
            fi
            COMPREPLY=( $(compgen -W "${allflags[*]}" -- "$cur") )
            if [[ $(type -t compopt) = "builtin" ]]; then
                [[ "${COMPREPLY[0]}" == *= ]] || compopt +o nospace
            fi

            # complete after --flag=abc
            if [[ $cur == *=* ]]; then
                if [[ $(type -t compopt) = "builtin" ]]; then
                    compopt +o nospace
                fi

                local index flag
                flag="${cur%=*}"
                __pi_index_of_word "${flag}" "${flags_with_completion[@]}"
                COMPREPLY=()
                if [[ ${index} -ge 0 ]]; then
                    PREFIX=""
                    cur="${cur#*=}"
                    ${flags_completion[${index}]}
                    if [ -n "${ZSH_VERSION}" ]; then
                        # zsh completion needs --flag= prefix
                        eval "COMPREPLY=( \"\${COMPREPLY[@]/#/${flag}=}\" )"
                    fi
                fi
            fi
            return 0;
            ;;
    esac

    # check if we are handling a flag with special work handling
    local index
    __pi_index_of_word "${prev}" "${flags_with_completion[@]}"
    if [[ ${index} -ge 0 ]]; then
        ${flags_completion[${index}]}
        return
    fi

    # we are parsing a flag and don't have a special handler, no completion
    if [[ ${cur} != "${words[cword]}" ]]; then
        return
    fi

    local completions
    completions=("${commands[@]}")
    if [[ ${#must_have_one_noun[@]} -ne 0 ]]; then
        completions=("${must_have_one_noun[@]}")
    fi
    if [[ ${#must_have_one_flag[@]} -ne 0 ]]; then
        completions+=("${must_have_one_flag[@]}")
    fi
    COMPREPLY=( $(compgen -W "${completions[*]}" -- "$cur") )

    if [[ ${#COMPREPLY[@]} -eq 0 && ${#noun_aliases[@]} -gt 0 && ${#must_have_one_noun[@]} -ne 0 ]]; then
        COMPREPLY=( $(compgen -W "${noun_aliases[*]}" -- "$cur") )
    fi

    if [[ ${#COMPREPLY[@]} -eq 0 ]]; then
        declare -F __custom_func >/dev/null && __custom_func
    fi

    # available in bash-completion >= 2, not always present on macOS
    if declare -F __ltrim_colon_completions >/dev/null; then
        __ltrim_colon_completions "$cur"
    fi

    # If there is only 1 completion and it is a flag with an = it will be completed
    # but we don't want a space after the =
    if [[ "${#COMPREPLY[@]}" -eq "1" ]] && [[ $(type -t compopt) = "builtin" ]] && [[ "${COMPREPLY[0]}" == --*= ]]; then
       compopt -o nospace
    fi
}

# The arguments should be in the form "ext1|ext2|extn"
__pi_handle_filename_extension_flag()
{
    local ext="$1"
    _filedir "@(${ext})"
}

__pi_handle_subdirs_in_dir_flag()
{
    local dir="$1"
    pushd "${dir}" >/dev/null 2>&1 && _filedir -d && popd >/dev/null 2>&1
}

__pi_handle_flag()
{
    __pi_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    # if a command required a flag, and we found it, unset must_have_one_flag()
    local flagname=${words[c]}
    local flagvalue
    # if the word contained an =
    if [[ ${words[c]} == *"="* ]]; then
        flagvalue=${flagname#*=} # take in as flagvalue after the =
        flagname=${flagname%=*} # strip everything after the =
        flagname="${flagname}=" # but put the = back
    fi
    __pi_debug "${FUNCNAME[0]}: looking for ${flagname}"
    if __pi_contains_word "${flagname}" "${must_have_one_flag[@]}"; then
        must_have_one_flag=()
    fi

    # if you set a flag which only applies to this command, don't show subcommands
    if __pi_contains_word "${flagname}" "${local_nonpersistent_flags[@]}"; then
      commands=()
    fi

    # keep flag value with flagname as flaghash
    # flaghash variable is an associative array which is only supported in bash > 3.
    if [[ -z "${BASH_VERSION}" || "${BASH_VERSINFO[0]}" -gt 3 ]]; then
        if [ -n "${flagvalue}" ] ; then
            flaghash[${flagname}]=${flagvalue}
        elif [ -n "${words[ $((c+1)) ]}" ] ; then
            flaghash[${flagname}]=${words[ $((c+1)) ]}
        else
            flaghash[${flagname}]="true" # pad "true" for bool flag
        fi
    fi

    # skip the argument to a two word flag
    if __pi_contains_word "${words[c]}" "${two_word_flags[@]}"; then
        c=$((c+1))
        # if we are looking for a flags value, don't show commands
        if [[ $c -eq $cword ]]; then
            commands=()
        fi
    fi

    c=$((c+1))

}

__pi_handle_noun()
{
    __pi_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    if __pi_contains_word "${words[c]}" "${must_have_one_noun[@]}"; then
        must_have_one_noun=()
    elif __pi_contains_word "${words[c]}" "${noun_aliases[@]}"; then
        must_have_one_noun=()
    fi

    nouns+=("${words[c]}")
    c=$((c+1))
}

__pi_handle_command()
{
    __pi_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    local next_command
    if [[ -n ${last_command} ]]; then
        next_command="_${last_command}_${words[c]//:/__}"
    else
        if [[ $c -eq 0 ]]; then
            next_command="_pi_root_command"
        else
            next_command="_${words[c]//:/__}"
        fi
    fi
    c=$((c+1))
    __pi_debug "${FUNCNAME[0]}: looking for ${next_command}"
    declare -F "$next_command" >/dev/null && $next_command
}

__pi_handle_word()
{
    if [[ $c -ge $cword ]]; then
        __pi_handle_reply
        return
    fi
    __pi_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"
    if [[ "${words[c]}" == -* ]]; then
        __pi_handle_flag
    elif __pi_contains_word "${words[c]}" "${commands[@]}"; then
        __pi_handle_command
    elif [[ $c -eq 0 ]]; then
        __pi_handle_command
    else
        __pi_handle_noun
    fi
    __pi_handle_word
}

_pi_admin_health-check()
{
    last_command="pi_admin_health-check"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_admin_version()
{
    last_command="pi_admin_version"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_admin()
{
    last_command="pi_admin"
    commands=()
    commands+=("health-check")
    commands+=("version")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_configure()
{
    last_command="pi_configure"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--APIHost=")
    flags+=("--ClientID=")
    flags+=("--ClientSecret=")
    flags+=("--IssuerID=")
    flags+=("--TenantID=")
    flags+=("--Token=")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_dag_create()
{
    last_command="pi_dag_create"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dagDesc=")
    flags+=("--dagFileName=")
    flags+=("--dagFilePath=")
    flags+=("--dagFlowType=")
    flags+=("--dagName=")
    flags+=("--dagTemplate=")
    flags+=("--dagVersion=")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_dag_delete()
{
    last_command="pi_dag_delete"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dagName=")
    flags+=("--force")
    flags+=("-f")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_dag_deploy()
{
    last_command="pi_dag_deploy"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dagName=")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_dag_list()
{
    last_command="pi_dag_list"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dagName=")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_dag_list-run()
{
    last_command="pi_dag_list-run"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dagName=")
    flags+=("--dagRunID=")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_dag_list-task()
{
    last_command="pi_dag_list-task"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dagName=")
    flags+=("--dagTaskID=")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_dag_status()
{
    last_command="pi_dag_status"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dagName=")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_dag_task-run-info()
{
    last_command="pi_dag_task-run-info"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dagName=")
    flags+=("--dagRunID=")
    flags+=("--dagTaskID=")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_dag_update()
{
    last_command="pi_dag_update"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dagDesc=")
    flags+=("--dagFileName=")
    flags+=("--dagFilePath=")
    flags+=("--dagFlowType=")
    flags+=("--dagName=")
    flags+=("--dagTemplate=")
    flags+=("--dagVersion=")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_dag()
{
    last_command="pi_dag"
    commands=()
    commands+=("create")
    commands+=("delete")
    commands+=("deploy")
    commands+=("list")
    commands+=("list-run")
    commands+=("list-task")
    commands+=("status")
    commands+=("task-run-info")
    commands+=("update")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_dependency_create()
{
    last_command="pi_dependency_create"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dependencyFileLocation=")
    flags+=("--dependencyFileName=")
    flags+=("--dependencyType=")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_dependency_delete()
{
    last_command="pi_dependency_delete"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dependencyID=")
    flags+=("--force")
    flags+=("-f")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_dependency_deploy()
{
    last_command="pi_dependency_deploy"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dependencyID=")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_dependency_list()
{
    last_command="pi_dependency_list"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dependencyID=")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_dependency_undeploy()
{
    last_command="pi_dependency_undeploy"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dependencyID=")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_dependency()
{
    last_command="pi_dependency"
    commands=()
    commands+=("create")
    commands+=("delete")
    commands+=("deploy")
    commands+=("list")
    commands+=("undeploy")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_flow_add-config-file()
{
    last_command="pi_flow_add-config-file"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--configFileDetails=")
    flags+=("--flowID=")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_flow_create()
{
    last_command="pi_flow_create"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--flowName=")
    flags+=("--flowTemplateID=")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_flow_create-direct()
{
    last_command="pi_flow_create-direct"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--desc=")
    flags+=("--flowFileName=")
    flags+=("--flowFilePath=")
    flags+=("--flowName=")
    flags+=("--flowType=")
    flags+=("--flowVersion=")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_flow_create-flow-template()
{
    last_command="pi_flow_create-flow-template"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--flowID=")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_flow_delete()
{
    last_command="pi_flow_delete"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--flowID=")
    flags+=("--force")
    flags+=("-f")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_flow_delete-config-file()
{
    last_command="pi_flow_delete-config-file"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--configFileName=")
    flags+=("--flowID=")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_flow_launch()
{
    last_command="pi_flow_launch"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--flowID=")
    flags+=("--flowTemplateID=")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_flow_list()
{
    last_command="pi_flow_list"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--flowID=")
    flags+=("--flowName=")
    flags+=("--flowTemplateID=")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_flow_list-config-files()
{
    last_command="pi_flow_list-config-files"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--flowID=")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_flow_list-tags()
{
    last_command="pi_flow_list-tags"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--flowID=")
    flags+=("--flowTemplateID=")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_flow_save-tags()
{
    last_command="pi_flow_save-tags"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--flowID=")
    flags+=("--flowTemplateID=")
    flags+=("--tags=")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_flow_stop()
{
    last_command="pi_flow_stop"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--flowName=")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_flow_update-direct()
{
    last_command="pi_flow_update-direct"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--desc=")
    flags+=("--flowFileName=")
    flags+=("--flowFilePath=")
    flags+=("--flowID=")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_flow_update-spark-args()
{
    last_command="pi_flow_update-spark-args"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--flowID=")
    flags+=("--flowTemplateID=")
    flags+=("--sparkArgs=")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_flow()
{
    last_command="pi_flow"
    commands=()
    commands+=("add-config-file")
    commands+=("create")
    commands+=("create-direct")
    commands+=("create-flow-template")
    commands+=("delete")
    commands+=("delete-config-file")
    commands+=("launch")
    commands+=("list")
    commands+=("list-config-files")
    commands+=("list-tags")
    commands+=("save-tags")
    commands+=("stop")
    commands+=("update-direct")
    commands+=("update-spark-args")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_flow-template_create()
{
    last_command="pi_flow-template_create"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--desc=")
    flags+=("--flowTemplateName=")
    flags+=("--flowTemplateVersion=")
    flags+=("--flowType=")
    flags+=("--templateFileName=")
    flags+=("--templateFilePath=")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_flow-template_delete()
{
    last_command="pi_flow-template_delete"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--flowTemplateID=")
    flags+=("--force")
    flags+=("-f")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_flow-template_list()
{
    last_command="pi_flow-template_list"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--flowTemplateID=")
    flags+=("--flowTemplateName=")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_flow-template_list-tags()
{
    last_command="pi_flow-template_list-tags"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--flowTemplateID=")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_flow-template_save-tags()
{
    last_command="pi_flow-template_save-tags"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--flowTemplateID=")
    flags+=("--tags=")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_flow-template_update()
{
    last_command="pi_flow-template_update"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--desc=")
    flags+=("--flowTemplateID=")
    flags+=("--flowTemplateName=")
    flags+=("--flowTemplateVersion=")
    flags+=("--flowType=")
    flags+=("--templateFileName=")
    flags+=("--templateFilePath=")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_flow-template_update-spark-args()
{
    last_command="pi_flow-template_update-spark-args"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--flowTemplateID=")
    flags+=("--sparkArgs=")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_flow-template()
{
    last_command="pi_flow-template"
    commands=()
    commands+=("create")
    commands+=("delete")
    commands+=("list")
    commands+=("list-tags")
    commands+=("save-tags")
    commands+=("update")
    commands+=("update-spark-args")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_instance_list-app-stages()
{
    last_command="pi_instance_list-app-stages"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--attemptID=")
    flags+=("--instanceID=")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_instance_list-attempt-details()
{
    last_command="pi_instance_list-attempt-details"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--attemptID=")
    flags+=("--instanceID=")
    flags+=("--stageAttemptID=")
    flags+=("--stageID=")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_instance_list-attempts()
{
    last_command="pi_instance_list-attempts"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--attemptID=")
    flags+=("--instanceID=")
    flags+=("--stageID=")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_instance_list-container-logs()
{
    last_command="pi_instance_list-container-logs"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--containerID=")
    flags+=("--containerLogSink=")
    flags+=("--instanceID=")
    flags+=("--tail")
    flags+=("-t")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_instance_list-container-response()
{
    last_command="pi_instance_list-container-response"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--containerID=")
    flags+=("--instanceID=")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_instance_list-containers()
{
    last_command="pi_instance_list-containers"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--instanceID=")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_instance_list-instance()
{
    last_command="pi_instance_list-instance"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--instanceID=")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_instance_list-spark-app-details()
{
    last_command="pi_instance_list-spark-app-details"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--instanceID=")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_instance_list-spark-executor-details()
{
    last_command="pi_instance_list-spark-executor-details"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--attemptID=")
    flags+=("--instanceID=")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_instance_list-submit-logs()
{
    last_command="pi_instance_list-submit-logs"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--instanceID=")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_instance_list-tasks()
{
    last_command="pi_instance_list-tasks"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--attemptID=")
    flags+=("--instanceID=")
    flags+=("--stageAttemptID=")
    flags+=("--stageID=")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_instance_stop()
{
    last_command="pi_instance_stop"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--instanceID=")
    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_instance()
{
    last_command="pi_instance"
    commands=()
    commands+=("list-app-stages")
    commands+=("list-attempt-details")
    commands+=("list-attempts")
    commands+=("list-container-logs")
    commands+=("list-container-response")
    commands+=("list-containers")
    commands+=("list-instance")
    commands+=("list-spark-app-details")
    commands+=("list-spark-executor-details")
    commands+=("list-submit-logs")
    commands+=("list-tasks")
    commands+=("stop")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_pi_root_command()
{
    last_command="pi"
    commands=()
    commands+=("admin")
    commands+=("configure")
    commands+=("dag")
    commands+=("dependency")
    commands+=("flow")
    commands+=("flow-template")
    commands+=("instance")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--config=")
    flags+=("--interactive")
    flags+=("-i")
    flags+=("--verbose")
    flags+=("-v")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

__start_pi()
{
    local cur prev words cword
    declare -A flaghash 2>/dev/null || :
    if declare -F _init_completion >/dev/null 2>&1; then
        _init_completion -s || return
    else
        __pi_init_completion -n "=" || return
    fi

    local c=0
    local flags=()
    local two_word_flags=()
    local local_nonpersistent_flags=()
    local flags_with_completion=()
    local flags_completion=()
    local commands=("pi")
    local must_have_one_flag=()
    local must_have_one_noun=()
    local last_command
    local nouns=()

    __pi_handle_word
}

if [[ $(type -t compopt) = "builtin" ]]; then
    complete -o default -F __start_pi pi
else
    complete -o default -o nospace -F __start_pi pi
fi

# ex: ts=4 sw=4 et filetype=sh
