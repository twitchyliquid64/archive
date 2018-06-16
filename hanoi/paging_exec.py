import paging_store, rules_store
import conf
import time
import smtplib
from email.mime.multipart import MIMEMultipart
from email.mime.text import MIMEText

last_paging_rule_run = 0 # epoch seconds

def should_run_paging_rules():
    global last_paging_rule_run
    if (last_paging_rule_run+80) < int(time.time()):
        last_paging_rule_run = int(time.time())
        return True
    return False


def send_email(user, pwd, recipient, subject, body):
    # Prepare actual message
    msg = MIMEMultipart('alternative')
    msg['Subject'] = subject
    msg['From'] = user
    msg['To'] = recipient
    htmlcontent = MIMEText(body, 'html')
    msg.attach(htmlcontent)

    try:
        server = smtplib.SMTP("smtp.gmail.com", 587)
        server.ehlo()
        server.starttls()
        server.login(user, pwd)
        server.sendmail(user, recipient, msg.as_string())
        server.close()
        print 'successfully sent the mail'
    except Exception, e:
        import traceback
        traceback.print_exc()
        print "failed to send mail", e

def generate_page_body(rule, state, rule_states):
    output = "You are being paged because this address is specified for the  <b>" + state.name + "</b> alert. This alert is failing in the <b>" + state.group + "</b> monitoring group.<br>\n"
    output += "The following rules are failing:<br><ul>\n"
    didfind = False
    for rule in rule_states:
        if rule_states[rule] == False:
            didfind = True
            output += "<li>" + rule + "</li>\n"
    if not didfind:
        output += "<p>N/A</p>"

    output += "</ul>\n<br>The following rules are succeeding:<br><ul>\n"
    didfind = False
    for rule in rule_states:
        if rule_states[rule] == True:
            didfind = True
            output += "<li>" + rule + "</li>\n"
    if not didfind:
        output += "<p>N/A</p>"

    output += "\n</ul><br>Please take appropriate action.<br>\n<br><i>-Hanoi</i>"
    return output

def page(rule, state, rule_states):
    if 'page' in rule:
        page_type = rule.get('page').get('type')
        if page_type == "GMAIL":
            subject = "[HANOI ALERT] Paging for " + state.name
            if 'subject' in rule.get('page'):
                subject = rule.get('page')['subject']
            body = generate_page_body(rule, state, rule_states)
            send_email(conf.get_config()['gmail']['username'], conf.get_config()['gmail']['password'], rule.get('page').get('address'), subject, body)


def exec_paging_rules():
    config = conf.get_config()
    for group in config['groups']:
        rules = config['groups'][group]['paging_rules']
        for rule in rules:
            state = paging_store.get(rule['name'], group)
            print "Running paging rule: %s.%s" % (state.group, state.name)
            hysteresis = rule.get('hysteresis', 180) * 60 #180 mins default

            new_state = {}
            has_found_failure = False
            for exec_rule_name in rule['rules']:
                ruleState = rules_store.getIfExists(exec_rule_name, group)
                if not ruleState or ruleState.ok:
                    new_state[exec_rule_name] = True
                else:
                    new_state[exec_rule_name] = False
                    has_found_failure = True

            if not state.failing and has_found_failure:
                if (state.last_page_epoch_seconds+hysteresis) < int(time.time()):
                    state.last_page_epoch_seconds = int(time.time())
                    print "Paging for: %s.%s" % (state.group, state.name)
                    page(rule, state, new_state)
                else:
                    print "Would page except hysteresis"
            state.rules_by_state = new_state
            state.failing = has_found_failure
